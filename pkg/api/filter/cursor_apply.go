package filter

import (
	"fmt"

	"gorm.io/gorm"
)

// apply cursor to filter
func ApplyCursor(
	query *gorm.DB,
	req Request,
	cfg Config,
) (*gorm.DB, error) {

	limit := req.Limit

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

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

	dir := "ASC"
	op := ">"

	if req.Desc {
		dir = "DESC"
		op = "<"
	}

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
