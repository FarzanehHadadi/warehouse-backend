package filter

var ThresholdProximityFilterConfig = FilterConfig{
	Filters: []FieldFilterConfig{
		{QueryParam: "product_id", Field: "product_id", Type: UIntType, Operator: Eq},
		{QueryParam: "store_id", Field: "store_id", Type: UIntType, Operator: Eq},
		{QueryParam: "category_id", Field: "category_id", Type: UIntType, Operator: Eq},
	},
	SortableFields: map[string]bool{
		"current_quantity": true,
		"product_id":       true,
	},
}

var StoreProductQuantityFilterConfig = FilterConfig{
	Filters: []FieldFilterConfig{
		{QueryParam: "product_id", Field: "product_id", Type: UIntType, Operator: Eq},
		{QueryParam: "category_id", Field: "category_id", Type: UIntType, Operator: Eq},
	},
	SortableFields: map[string]bool{
		"total_quantity": true,
		"product_id":     true,
		"store_id":       true,
	},
}
