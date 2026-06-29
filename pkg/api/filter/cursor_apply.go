package filter

import (
	"fmt"

	"gorm.io/gorm"
)

// CompositeKeysetCondition builds a WHERE clause for (primaryExpr, tieBreaker) keyset pagination.
func CompositeKeysetCondition(primaryExpr, tieBreaker string, desc bool, primaryValue any, tieValue uint) (string, []any) {
	op := KeysetOp(desc)
	return fmt.Sprintf("(%s, %s) %s (?, ?)", primaryExpr, tieBreaker, op), []any{primaryValue, tieValue}
}

// SingleKeysetCondition builds a WHERE clause for single-column keyset pagination.
func SingleKeysetCondition(field string, desc bool, value any) (string, []any) {
	op := KeysetOp(desc)
	return fmt.Sprintf("%s %s ?", field, op), []any{value}
}

// TripleKeysetCondition builds a WHERE clause for three-column keyset pagination.
func TripleKeysetCondition(expr1, expr2, expr3 string, desc bool, val1, val2 any, val3 uint) (string, []any) {
	op := KeysetOp(desc)
	return fmt.Sprintf("(%s, %s, %s) %s (?, ?, ?)", expr1, expr2, expr3, op), []any{val1, val2, val3}
}

// ApplyCursor applies GORM cursor pagination for models sorted by created_at.
func ApplyCursor(
	query *gorm.DB,
	req Request,
	cfg FilterConfig,
) (*gorm.DB, error) {

	limit := NormalizeLimit(req.Limit)

	sortField := "created_at"

	if req.SortBy != "" {
		if !cfg.SortableFields[req.SortBy] {
			return nil, fmt.Errorf(
				"field %s is not sortable",
				req.SortBy,
			)
		}

		sortField = req.SortBy
	}

	dir := SortSQLDir(req.Desc)
	op := KeysetOp(req.Desc)

	cursor, err := DecodeCursor(req.Cursor)
	if err != nil {
		return nil, err
	}

	if cursor != nil {

		if sortField != "created_at" {
			return nil, fmt.Errorf(
				"cursor pagination currently supports only created_at sorting",
			)
		}

		query = query.Where(
			fmt.Sprintf("(created_at, id) %s (?, ?)", op),
			cursor.CreatedAt,
			cursor.ID,
		)
	}

	query = query.
		Order(fmt.Sprintf("%s %s, id %s",
			sortField,
			dir,
			dir,
		)).
		Limit(limit + 1)

	return query, nil
}
