package filter

var ProductFilterConfig = Config{
	Filters: []FieldFilterConfig{
		{QueryParam: "name", Field: "name", Type: StringType, Operator: Like},
		{QueryParam: "created_after", Field: "created_at", Type: TimeType, Operator: Gt},
		{QueryParam: "created_from", Field: "created_at", Type: TimeType, Operator: Gte},
	},
	SearchFields: []string{"name"},
	SortableFields: map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	},
}
