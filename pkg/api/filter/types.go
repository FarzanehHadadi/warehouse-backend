package filter

type Operator string

// define operators type
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

// define type string,uint,time
type FieldType string

const (
	StringType  FieldType = "string"
	UIntType    FieldType = "uint"
	TimeType    FieldType = "time"
	BooleanType FieldType = "boolean"
)

// define condition
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

type Config struct {
	Fields         map[string]FieldType
	SearchFields   []string
	SortableFields map[string]bool
}

type CursorResponse struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}
