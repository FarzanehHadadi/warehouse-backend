package filter

// NormalizeLimit clamps list page size to the API default range.
func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// ResolveSort returns the effective sort field when it is allowed by cfg, otherwise defaultField.
func ResolveSort(req Request, cfg FilterConfig, defaultField string) string {
	if req.SortBy != "" && cfg.SortableFields[req.SortBy] {
		return req.SortBy
	}
	return defaultField
}

// SortSQLDir returns "ASC" or "DESC" for SQL ORDER BY clauses.
func SortSQLDir(desc bool) string {
	if desc {
		return "DESC"
	}
	return "ASC"
}

// KeysetOp returns ">" for ascending keyset pagination and "<" for descending.
func KeysetOp(desc bool) string {
	if desc {
		return "<"
	}
	return ">"
}

// UintFilterValue returns the uint value of an equality filter, if present.
func UintFilterValue(req Request, field string) (uint, bool) {
	for _, cond := range req.Filters {
		if cond.Field != field || cond.Operator != Eq || len(cond.Values) == 0 {
			continue
		}
		v, ok := cond.Values[0].(uint)
		return v, ok
	}
	return 0, false
}
