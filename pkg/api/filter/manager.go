package filter

var ManagerFilterConfig = FilterConfig{
	Filters: []FieldFilterConfig{
		{QueryParam: "id", Field: "id", Type: UIntType, Operator: Eq},
		{QueryParam: "name", Field: "name", Type: StringType, Operator: Like},
		{QueryParam: "phone", Field: "phone", Type: StringType, Operator: Like},
		{QueryParam: "created_after", Field: "created_at", Type: TimeType, Operator: Gt},
		{QueryParam: "created_before", Field: "created_at", Type: TimeType, Operator: Lte},
		{QueryParam: "status", Field: "status", Type: StringType, Operator: In},
	},
	SearchFields: []string{"name", "phone"},
	SortableFields: map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	},
}
