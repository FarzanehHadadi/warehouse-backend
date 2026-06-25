package filter

type Operator string

const (
	Eq   Operator = "eq"
	Ne   Operator = "ne"
	Gt   Operator = "gt"
	Gte  Operator = "gte"
	Lt   Operator = "lt"
	Lte  Operator = "lte"
	Like Operator = "like"
	In   Operator = "in"
)

type FieldType string

const (
	StringType  FieldType = "string"
	UIntType    FieldType = "uint"
	TimeType    FieldType = "time"
	BooleanType FieldType = "boolean"
)

type Condition struct {
	Field    string   `json:"field"`
	Operator Operator `json:"operator"`
	Values   []any    `json:"values"`
}

type Request struct {
	Filters []Condition `json:"filters"`
	Search  string      `json:"search"`
	SortBy  string      `json:"sort_by"`
	Desc    bool        `json:"desc"`
	Cursor  string      `json:"cursor"`
	Limit   int         `json:"limit"`
}

type FieldFilterConfig struct {
	QueryParam string    // URL query key, e.g. "name", "created_after", "status"
	Field      string    // DB column used in WHERE
	Type       FieldType // string | uint | time | boolean
	Operator   Operator  // eq | ne | gt | gte | lt | lte | like | in
}

type FilterConfig struct {
	Filters        []FieldFilterConfig
	SearchFields   []string
	SortableFields map[string]bool
}

func (cfg FilterConfig) FieldTypes() map[string]FieldType {
	types := make(map[string]FieldType, len(cfg.Filters))
	for _, f := range cfg.Filters {
		types[f.Field] = f.Type
	}
	return types
}

type CursorResponse struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}
