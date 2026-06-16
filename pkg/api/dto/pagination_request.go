package dto

import (
	"strconv"
	"strings"
	"time"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var reservedQueryParams = map[string]struct{}{
	"limit":      {},
	"sort_by":    {},
	"sort_order": {},
	"search":     {},
	"cursor":     {},
}

func NewPaginationRequestFromConfig(c *gin.Context, cfg filter.Config) *filter.Request {
	req := &filter.Request{
		Limit:  20,
		SortBy: "created_at",
		Desc:   true,
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	req.SortBy = c.DefaultQuery("sort_by", "created_at")
	if !cfg.SortableFields[req.SortBy] {
		req.SortBy = "created_at"
	}

	req.Desc = strings.ToLower(c.DefaultQuery("sort_order", "desc")) == "desc"
	req.Search = c.Query("search")
	req.Cursor = c.Query("cursor")

	req = parseFieldFiltersFromConfig(c, req, cfg)
	return Normalize(req)
}

func Normalize(r *filter.Request) *filter.Request {
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.Limit > 100 {
		r.Limit = 100
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
	return r
}

func parseFieldFiltersFromConfig(c *gin.Context, r *filter.Request, cfg filter.Config) *filter.Request {
	for _, fc := range cfg.Filters {

		if _, reserved := reservedQueryParams[fc.QueryParam]; reserved {
			continue
		}

		rawValues := readQueryValues(c, fc.QueryParam, fc.Operator == filter.In)
		logger.Log.Info("processing filter",
			zap.Any("rawValues", rawValues),
		)
		logger.Log.Info("processing filter",
			zap.Any("fc", fc),
		)
		if len(rawValues) == 0 {
			continue
		}

		if cond, ok := buildConditionFromFieldConfig(fc, rawValues); ok {
			r.Filters = append(r.Filters, cond)
		}
	}
	return r
}

func readQueryValues(c *gin.Context, key string, allowMultiple bool) []string {
	if allowMultiple {
		if values := c.QueryArray(key); len(values) > 0 {
			return expandCommaSeparated(values)
		}
	}

	if v := c.Query(key); v != "" {
		if allowMultiple {
			return expandCommaSeparated([]string{v})
		}
		return []string{v}
	}

	return nil
}

func expandCommaSeparated(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}

func buildConditionFromFieldConfig(fc filter.FieldFilterConfig, rawValues []string) (filter.Condition, bool) {
	parsed, ok := parseFilterValues(fc.Type, fc.Operator, rawValues)
	if !ok || len(parsed) == 0 {
		return filter.Condition{}, false
	}

	return filter.Condition{
		Field:    fc.Field,
		Operator: fc.Operator,
		Values:   parsed,
	}, true
}

func parseFilterValues(fieldType filter.FieldType, op filter.Operator, rawValues []string) ([]any, bool) {
	if op == filter.In {
		return parseMany(fieldType, rawValues)
	}
	if len(rawValues) != 1 {
		return nil, false
	}
	value, ok := parseOne(fieldType, rawValues[0])
	if !ok {
		return nil, false
	}
	return []any{value}, true
}

func parseMany(fieldType filter.FieldType, rawValues []string) ([]any, bool) {
	out := make([]any, 0, len(rawValues))
	for _, raw := range rawValues {
		value, ok := parseOne(fieldType, raw)
		if !ok {
			return nil, false
		}
		out = append(out, value)
	}
	return out, true
}

func parseOne(fieldType filter.FieldType, raw string) (any, bool) {
	switch fieldType {
	case filter.StringType:
		return raw, true
	case filter.UIntType:
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, false
		}
		return uint(n), true
	case filter.BooleanType:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, false
		}
		return b, true
	case filter.TimeType:
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			return t, true
		}
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return nil, false
		}
		return t, true
	default:
		return nil, false
	}
}
func NewPagination(c *gin.Context, cfg filter.Config) *filter.Request {
	filters := &filter.Request{
		Limit:  20,
		SortBy: "created_at",
		Desc:   true,
	}

	limit := c.Query("limit")
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	if cfg.SortableFields[sortBy] {
		filters.SortBy = sortBy
	}

	filters.Desc = strings.ToLower(c.DefaultQuery("sort_order", "desc")) == "desc"
	filters.Search = c.Query("search")
	filters.Cursor = c.Query("cursor")

	filters = parseFilter(c, filters, cfg)
	return normalize(filters)
}
func normalize(req *filter.Request) *filter.Request {
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	return req
}
func parseFilter(c *gin.Context, req *filter.Request, cfg filter.Config) *filter.Request {
	return req
}
