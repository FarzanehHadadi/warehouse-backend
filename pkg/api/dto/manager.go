package dto

type DepartmentSummary struct {
	ID   uint   `json:"id" example:"1"`             // Department ID
	Name string `json:"name" example:"Engineering"` // Department name
}
type ManagerSummary struct {
	ID          uint                `json:"id" example:"1"`
	Name        string              `json:"name" example:"John Doe"`
	Phone       string              `json:"phone" example:"+1234567890"`
	Departments []DepartmentSummary `json:"departments"`
}

// ManagerListResponse represents a paginated list of managers
type ManagerListResponse struct {
	Items []ManagerSummary `json:"items"`
	Total int64            `json:"total" example:"10"`
	// Page  int              `json:"page" example:"1"`
	// Limit int              `json:"limit" example:"20"`
}

// CreateManagerRequest represents request body for creating a manager
type CreateManagerRequest struct {
	Name  string `json:"name" binding:"required,min=3,max=100" example:"John Doe"`
	Phone string `json:"phone" binding:"required,min=7,max=11" example:"+1234567890"`
}

// UpdateManagerRequest represents request body for updating a manager
type UpdateManagerRequest struct {
	Name  string `json:"name" binding:"omitempty,min=3,max=100" example:"John Doe Updated"`
	Phone string `json:"phone" binding:"omitempty,min=7,max=11" example:"+1987654321"`
}
