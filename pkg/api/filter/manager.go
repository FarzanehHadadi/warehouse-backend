package filter

var ManagerFilterConfig = Config{
	Filters: []FieldFilterConfig{
		{QueryParam: "id", Field: "id", Type: UIntType, Operator: Eq},
		{QueryParam: "name", Field: "name", Type: StringType, Operator: Like},
		{QueryParam: "phone", Field: "phone", Type: StringType, Operator: Like},
		{QueryParam: "created_after", Field: "created_at", Type: TimeType, Operator: Gt},
		{QueryParam: "created_from", Field: "created_at", Type: TimeType, Operator: Gte},
		{QueryParam: "status", Field: "status", Type: StringType, Operator: In},
	},
	SearchFields: []string{"name", "phone"},
	SortableFields: map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	},
}
