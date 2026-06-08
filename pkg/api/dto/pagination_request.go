package dto

import (
	"strconv"
	"strings"
	"time"
	"warehouse/pkg/api/filter"

	"github.com/gin-gonic/gin"
)

// NewPaginationRequest - Clean helper you liked
func NewPaginationRequest(c *gin.Context) *filter.Request {
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
	req.Desc = strings.ToLower(c.DefaultQuery("sort_order", "desc")) == "desc"
	req.Search = c.Query("search")
	req.Cursor = c.Query("cursor")
	req = parseFieldFilters(c, req)
	req = Normalize(req)
	return req
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
func parseFieldFilters(c *gin.Context, r *filter.Request) *filter.Request {
	// Name - LIKE
	if name := c.Query("name"); name != "" {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "name",
			Operator: filter.Like,
			Values:   []any{name},
		})
	}

	// Phone - LIKE (if you want)
	if phone := c.Query("phone"); phone != "" {
		r.Filters = append(r.Filters, filter.Condition{
			Field:    "phone",
			Operator: filter.Like,
			Values:   []any{phone},
		})
	}

	// Created At >
	if gt := c.Query("created_at_gt"); gt != "" {
		if t, err := time.Parse("2006-01-02", gt); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "created_at",
				Operator: filter.Gt,
				Values:   []any{t},
			})
		}
	}

	// Created At >=
	if gte := c.Query("created_at_gte"); gte != "" {
		if t, err := time.Parse("2006-01-02", gte); err == nil {
			r.Filters = append(r.Filters, filter.Condition{
				Field:    "created_at",
				Operator: filter.Gte,
				Values:   []any{t},
			})
		}
	}
	return r

}
