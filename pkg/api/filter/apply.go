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

	return query, nil
}

func applyString(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq:
		return query.Where(cond.Field+" = ?", cond.Values[0])
	case Ne:
		return query.Where(cond.Field+" != ?", cond.Values[0])
	case Like:
		return query.Where(cond.Field+" ILIKE ?", "%"+fmt.Sprint(cond.Values[0])+"%")
	case In:
		return query.Where(cond.Field+" IN ?", cond.Values)
	}
	return query
}

func applyUint(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq:
		return query.Where(cond.Field+" = ?", cond.Values[0])
	case Ne:
		return query.Where(cond.Field+" != ?", cond.Values[0])
	case Gt:
		return query.Where(cond.Field+" > ?", cond.Values[0])
	case Gte:
		return query.Where(cond.Field+" >= ?", cond.Values[0])
	case Lt:
		return query.Where(cond.Field+" < ?", cond.Values[0])
	case Lte:
		return query.Where(cond.Field+" <= ?", cond.Values[0])
	case In:
		return query.Where(cond.Field+" IN ?", cond.Values)
	}
	return query
}

func applyTime(query *gorm.DB, cond Condition) *gorm.DB {
	if len(cond.Values) == 0 {
		return query
	}

	switch cond.Operator {
	case Eq, Gt, Gte, Lt, Lte:
		return query.Where(cond.Field+" "+getSQLOp(cond.Operator)+" ?", cond.Values[0])
	case In:
		return query.Where(cond.Field+" IN ?", cond.Values)
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
		return query.Where(cond.Field+" = ?", cond.Values[0])
	}
	return query
}
