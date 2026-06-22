package filter

var SimpleFilterConfig = Config{
	Filters: []FieldFilterConfig{
		{QueryParam: "name", Field: "name", Type: StringType, Operator: Like},
		{QueryParam: "created_after", Field: "created_at", Type: TimeType, Operator: Gt},
		{QueryParam: "created_before", Field: "created_at", Type: TimeType, Operator: Lte},
	},
	SearchFields: []string{"name"},
	SortableFields: map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	},
}
