package repository

import (
	"context"
	"fmt"

	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

const orderQuantityExpr = `
	SUM(
		CASE
			WHEN o.type = 'inbound' THEN o.quantity
			WHEN o.type = 'outbound' THEN -o.quantity
			ELSE 0
		END
	)
`

const thresholdProximityBaseQuery = `
SELECT
	p.id AS product_id,
	p.name AS product_name,
	p.category_id,
	c.name AS category_name,
	p.unit_id,
	u.name AS unit_name,
	p.warning_threshold,
	COALESCE(inv.current_quantity, 0) AS current_quantity,
	COALESCE(inv.stores_count, 0) AS stores_count
FROM products p
INNER JOIN categories c ON c.id = p.category_id
INNER JOIN units u ON u.id = p.unit_id
LEFT JOIN (
	SELECT
		per_store.product_id,
		COALESCE(SUM(GREATEST(per_store.quantity, 0)), 0) AS current_quantity,
		COUNT(*) FILTER (WHERE per_store.quantity > 0) AS stores_count
	FROM (
		SELECT
			o.product_id,
			o.store_id,
			SUM(
				CASE
					WHEN o.type = 'inbound' THEN o.quantity
					WHEN o.type = 'outbound' THEN -o.quantity
					ELSE 0
				END
			) AS quantity
		FROM orders o
		WHERE o.product_status = 'good'
		%s
		GROUP BY o.product_id, o.store_id
	) per_store
	GROUP BY per_store.product_id
) inv ON inv.product_id = p.id
WHERE COALESCE(inv.current_quantity, 0) <= p.warning_threshold
`

func (r *reportRepository) GetThresholdProximity(ctx context.Context, req filter.Request, isPaginated bool) ([]models.ThresholdProximityReport, *filter.CursorResponse, error) {
	cfg := filter.ThresholdProximityFilterConfig
	limit := filter.NormalizeLimit(req.Limit)

	args := make([]any, 0, 8)
	storeClause := ""
	if storeID, ok := filter.UintFilterValue(req, "store_id"); ok {
		storeClause = "AND o.store_id = ?"
		args = append(args, storeID)
	}

	query := fmt.Sprintf(thresholdProximityBaseQuery, storeClause)

	if productID, ok := filter.UintFilterValue(req, "product_id"); ok {
		query += " AND p.id = ?"
		args = append(args, productID)
	}
	if categoryID, ok := filter.UintFilterValue(req, "category_id"); ok {
		query += " AND p.category_id = ?"
		args = append(args, categoryID)
	}

	sortField := filter.ResolveSort(req, cfg, "current_quantity")
	sortDir := filter.SortSQLDir(req.Desc)

	if isPaginated {
		cursor, err := filter.DecodeReportCursor(req.Cursor)
		if err != nil {
			return nil, nil, err
		}
		if cursor != nil && cursor.ProductID > 0 {
			switch sortField {
			case "product_id":
				clause, clauseArgs := filter.SingleKeysetCondition("p.id", req.Desc, cursor.ProductID)
				query += " AND " + clause
				args = append(args, clauseArgs...)
			default:
				clause, clauseArgs := filter.CompositeKeysetCondition(
					"COALESCE(inv.current_quantity, 0)",
					"p.id",
					req.Desc,
					cursor.CurrentQuantity,
					cursor.ProductID,
				)
				query += " AND " + clause
				args = append(args, clauseArgs...)
			}
		}
	}

	switch sortField {
	case "product_id":
		query += fmt.Sprintf(" ORDER BY p.id %s", sortDir)
	default:
		query += fmt.Sprintf(" ORDER BY current_quantity %s, p.id %s", sortDir, sortDir)
	}

	if isPaginated {
		query += " LIMIT ?"
		args = append(args, limit+1)
	}

	var results []models.ThresholdProximityReport
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, nil, err
	}
	if results == nil {
		results = []models.ThresholdProximityReport{}
	}
	if !isPaginated {
		return results, nil, nil
	}

	results, cursorResp := filter.BuildPaginatedResult(results, limit, func(item models.ThresholdProximityReport) string {
		return filter.EncodeReportCursor(item.ProductID, item.CurrentQuantity)
	})

	return results, &cursorResp, nil
}

const storeProductQuantityBaseQuery = `
SELECT
	qty_report.product_id,
	qty_report.category_id,
	qty_report.store_id,
	qty_report.total_quantity
FROM (
	SELECT
		o.product_id,
		p.category_id,
		o.store_id,
` + orderQuantityExpr + ` AS total_quantity
	FROM orders o
	INNER JOIN products p ON p.id = o.product_id
	WHERE o.product_status = 'good'
`

func (r *reportRepository) GetStoreProductQuantities(ctx context.Context, req filter.Request, isPaginated bool) ([]models.StoreProductQuantityReport, *filter.CursorResponse, error) {
	cfg := filter.StoreProductQuantityFilterConfig
	limit := filter.NormalizeLimit(req.Limit)

	args := make([]any, 0, 8)
	query := storeProductQuantityBaseQuery

	if productID, ok := filter.UintFilterValue(req, "product_id"); ok {
		query += " AND o.product_id = ?"
		args = append(args, productID)
	}
	if categoryID, ok := filter.UintFilterValue(req, "category_id"); ok {
		query += " AND p.category_id = ?"
		args = append(args, categoryID)
	}

	query += " GROUP BY o.product_id, p.category_id, o.store_id) qty_report WHERE 1=1"

	sortField := filter.ResolveSort(req, cfg, "store_id")
	sortDir := filter.SortSQLDir(req.Desc)

	if isPaginated {
		cursor, err := filter.DecodeStoreProductQuantityCursor(req.Cursor)
		if err != nil {
			return nil, nil, err
		}
		if cursor != nil {
			switch sortField {
			case "total_quantity":
				clause, clauseArgs := filter.TripleKeysetCondition(
					"qty_report.total_quantity", "qty_report.store_id", "qty_report.product_id",
					req.Desc,
					cursor.TotalQuantity, cursor.StoreID, cursor.ProductID,
				)
				query += " AND " + clause
				args = append(args, clauseArgs...)
			case "product_id":
				clause, clauseArgs := filter.CompositeKeysetCondition(
					"qty_report.product_id", "qty_report.store_id",
					req.Desc,
					cursor.ProductID, cursor.StoreID,
				)
				query += " AND " + clause
				args = append(args, clauseArgs...)
			default:
				clause, clauseArgs := filter.CompositeKeysetCondition(
					"qty_report.store_id", "qty_report.product_id",
					req.Desc,
					cursor.StoreID, cursor.ProductID,
				)
				query += " AND " + clause
				args = append(args, clauseArgs...)
			}
		}
	}

	switch sortField {
	case "total_quantity":
		query += fmt.Sprintf(" ORDER BY qty_report.total_quantity %s, qty_report.store_id %s, qty_report.product_id %s", sortDir, sortDir, sortDir)
	case "product_id":
		query += fmt.Sprintf(" ORDER BY qty_report.product_id %s, qty_report.store_id %s", sortDir, sortDir)
	default:
		query += fmt.Sprintf(" ORDER BY qty_report.store_id %s, qty_report.product_id %s", sortDir, sortDir)
	}

	if isPaginated {
		query += " LIMIT ?"
		args = append(args, limit+1)
	}

	var results []models.StoreProductQuantityReport
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, nil, err
	}
	if results == nil {
		results = []models.StoreProductQuantityReport{}
	}
	if !isPaginated {
		return results, nil, nil
	}

	results, cursorResp := filter.BuildPaginatedResult(results, limit, func(item models.StoreProductQuantityReport) string {
		return filter.EncodeStoreProductQuantityCursor(item.TotalQuantity, item.StoreID, item.ProductID)
	})

	return results, &cursorResp, nil
}
