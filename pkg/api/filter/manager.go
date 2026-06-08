package filter

var ManagerFilterConfig = Config{
	Fields: map[string]FieldType{
		"id":         UIntType,
		"name":       StringType,
		"phone":      StringType,
		"created_at": TimeType,
		// "is_active":   BooleanType,
		// "email":       StringType,
		// "status":      StringType,
	},
	SearchFields: []string{"name", "phone"},
	SortableFields: map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
		// "status":     true,
	},
}
