package filter

var ProductFilterConfig = FilterConfig{
	Filters: []FieldFilterConfig{
		{QueryParam: "name", Field: "name", Type: StringType, Operator: Like},
		{QueryParam: "category_id", Field: "category_id", Type: UIntType, Operator: Eq},
		{QueryParam: "unit_id", Field: "unit_id", Type: UIntType, Operator: Eq},
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
