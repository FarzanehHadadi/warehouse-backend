package filter

var OrderFilterConfig = Config{
	Filters: []FieldFilterConfig{
		{QueryParam: "product_id", Field: "product_id", Type: UIntType, Operator: Eq},
		{QueryParam: "store_id", Field: "store_id", Type: UIntType, Operator: Eq},
		{QueryParam: "department_id", Field: "department_id", Type: UIntType, Operator: Eq},
		{QueryParam: "type", Field: "type", Type: StringType, Operator: In},
		{QueryParam: "product_status", Field: "product_status", Type: StringType, Operator: In},
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
