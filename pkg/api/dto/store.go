package dto

type StoreSummary struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Store 1"`
	ManagerName string `json:"manager_name" example:"John Doe"`
	IsActive    bool   `json:"is_active" example:"true"`
}
type storeListResponse struct {
	Items      []StoreSummary `json:"items"`
	NextCursor string         `json:"next_cursor" example:""`
	HasMore    bool           `json:"has_more" example:"true"`
	Limit      int            `json:"limit" example:"20"`
}
