// apply search and sort queries
package filter

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func Apply(query *gorm.DB, req Request, cfg Config) (*gorm.DB, error) {
	for _, cond := range req.Filters {
		fieldType, ok := cfg.Fields[cond.Field]
		if !ok {
			return nil, fmt.Errorf("field %s not allowed", cond.Field)
		}

		switch fieldType {
		case UIntType:
			query = applyUint(query, cond)
		case StringType:
			query = applyString(query, cond)
		case TimeType:
			query = applyTime(query, cond)
		}
	}

	if req.Search != "" && len(cfg.SearchFields) > 0 {
		query = applySearch(query, req.Search, cfg.SearchFields)
	}

	// Sorting
	if cfg.SortableFields[req.SortBy] {
		direction := "ASC"
		if req.Desc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, direction))
	} else if req.SortBy == "" {
		// Default sort
		query = query.Order("created_at DESC, id DESC")
	}

	return query, nil
}

func applyString(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq:
		return query.Where(gorm.Expr(cond.Field+" = ?"), cond.Values[0])
	case Ne:
		return query.Where(gorm.Expr(cond.Field+" != ?"), cond.Values[0])
	case Like:
		return query.Where(gorm.Expr(cond.Field+" ILIKE ?"), "%"+fmt.Sprint(cond.Values[0])+"%")
	case In:
		return query.Where(gorm.Expr(cond.Field+" IN ?"), cond.Values)
	}
	return query
}

func applyUint(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq:
		return query.Where(gorm.Expr(cond.Field+" = ?"), cond.Values[0])
	case Ne:
		return query.Where(gorm.Expr(cond.Field+" != ?"), cond.Values[0])
	case Gt:
		return query.Where(gorm.Expr(cond.Field+" > ?"), cond.Values[0])
	case Gte:
		return query.Where(gorm.Expr(cond.Field+" >= ?"), cond.Values[0])
	case Lt:
		return query.Where(gorm.Expr(cond.Field+" < ?"), cond.Values[0])
	case Lte:
		return query.Where(gorm.Expr(cond.Field+" <= ?"), cond.Values[0])
	case In:
		return query.Where(gorm.Expr(cond.Field+" IN ?"), cond.Values)
	}
	return query
}

func applyTime(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq, Gt, Gte, Lt, Lte:
		return query.Where(gorm.Expr(cond.Field+" "+getSQLOp(cond.Operator)+" ?"), cond.Values[0])
	case In:
		return query.Where(gorm.Expr(cond.Field+" IN ?"), cond.Values)
	}
	return query
}

func getSQLOp(op Operator) string {
	switch op {
	case Gt:
		return ">"
	case Gte:
		return ">="
	case Lt:
		return "<"
	case Lte:
		return "<="
	default:
		return "="
	}
}

func applySearch(query *gorm.DB, search string, fields []string) *gorm.DB {
	conditions := make([]string, len(fields))
	args := make([]any, len(fields))

	for i, field := range fields {
		conditions[i] = field + " ILIKE ?"
		args[i] = "%" + search + "%"
	}

	return query.Where(strings.Join(conditions, " OR "), args...)
}
func applyBoolean(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq:
		return query.Where(gorm.Expr(cond.Field+" = ?"), cond.Values[0])
	}
	return query
}

/*
package dto

import (
	"strconv"
	"strings"
	"time"
	"warehouse/pkg/filter"
)

type PaginationRequest struct {
	filter.Request
	Page int `json:"page,omitempty"`
}

func NewPaginationRequest(c *gin.Context) *PaginationRequest {
	req := &PaginationRequest{
		Request: filter.Request{
			Limit:   20,
			SortBy:  "created_at",
			Desc:    true,
		},
	}

	// Basic parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = l
		}
	}
	req.SortBy = c.DefaultQuery("sort_by", "created_at")
	req.SortOrder = c.DefaultQuery("sort_order", "desc")
	req.Desc = strings.ToLower(req.SortOrder) == "desc"
	req.Search = c.Query("search")
	req.Cursor = c.Query("cursor")

	// Parse all field filters
	req.parseFieldFilters(c)
	req.Normalize()

	return req
}

func (r *PaginationRequest) parseFieldFilters(c *gin.Context) {
	// 1. Name (Like)
	if name := c.Query("name"); name != "" {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "name",
			Operator: filter.Like,
			Values:   []any{name},
		})
	}

	// 2. Boolean (is_active)
	if isActive := c.Query("is_active"); isActive != "" {
		if val, err := strconv.ParseBool(isActive); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "is_active",
				Operator: filter.Eq,
				Values:   []any{val},
			})
		}
	}

	// 3. Status - Single OR Multiple
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		if len(statuses) == 1 {
			// Single status
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "status",
				Operator: filter.Eq,
				Values:   []any{statuses[0]},
			})
		} else {
			// Multiple statuses (IN operator)
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "status",
				Operator: filter.In,
				Values:   convertToAny(statuses),
			})
		}
	}

	// 4. Created At filters
	if gt := c.Query("created_at_gt"); gt != "" {
		if t, err := time.Parse("2006-01-02", gt); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "created_at",
				Operator: filter.Gt,
				Values:   []any{t},
			})
		}
	}

	// Add more filters as needed...
}

func convertToAny(slice []string) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}*/

/*
var ManagerFilterConfig = filter.Config{
    Fields: map[string]filter.FieldType{
        "id":          filter.UIntType,
        "name":        filter.StringType,
        "phone":       filter.StringType,
        "email":       filter.StringType,
        "created_at":  filter.TimeType,
        "is_active":   filter.BooleanType,
        "status":      filter.StringType,
    },
    SearchFields: []string{"name", "phone", "email"},
    SortableFields: map[string]bool{
        "id":         true,
        "name":       true,
        "created_at": true,
        "status":     true,
    },
}
*/
// @Summary      Get Managers List
// @Description  Retrieve managers with filtering, search, and cursor pagination
// @Tags         Managers
// @Accept       json
// @Produce      json
// @Param        search          query    string    false  "Global search in name, phone, email"
// @Param        name            query    string    false  "Filter by name (partial match)"
// @Param        is_active       query    boolean   false  "Filter by active status" Enums(true,false)
// @Param        status          query    string    false  "Filter by status. Use multiple times for multiple selection" CollectionFormat(multi) Enums(active,pending,inactive,suspended)
// @Param        created_at_gt   query    string    false  "Created after date (YYYY-MM-DD)"
// @Param        sort_by         query    string    false  "Sort field" Enums(id,name,created_at,status)
// @Param        sort_order      query    string    false  "Sort direction" Enums(asc,desc)
// @Param        cursor          query    string    false  "Cursor for next page"
// @Param        limit           query    integer   false  "Number of items per page (max 100)" minimum(1) maximum(100)
// @Success      200  {object}  map[string]interface{}  "data, next_cursor, has_more"
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/managers [get]

/*
// Add these new parsing methods
func (r *PaginationRequest) parseFieldFilters(c *gin.Context) {
	// 1. Boolean filter: is_active=true / is_active=false
	if isActive := c.Query("is_active"); isActive != "" {
		if val, err := strconv.ParseBool(isActive); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "is_active",
				Operator: filter.Eq,
				Values:   []any{val},
			})
		}
	}

	// 2. Status / Enum filter (single value)
	if status := c.Query("status"); status != "" {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "status",
			Operator: filter.Eq,
			Values:   []any{status},
		})
	}

	// 3. Status with IN operator (multiple values)
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "status",
			Operator: filter.In,
			Values:   convertToAny(statuses),
		})
	}

	// Existing filters...
	if name := c.Query("name"); name != "" {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "name",
			Operator: filter.Like,
			Values:   []any{name},
		})
	}

	if createdAtGt := c.Query("created_at_gt"); createdAtGt != "" {
		if t, err := time.Parse("2006-01-02", createdAtGt); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "created_at",
				Operator: filter.Gt,
				Values:   []any{t},
			})
		}
	}
}

// Helper function
func convertToAny(slice []string) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
*/
