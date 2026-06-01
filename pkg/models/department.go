package models

type Department struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required,min=1,max=100" gorm:"uniqueIndex; not null;size:100"`
	ManagerName string `json:"manager_name" binding:"omitempty,min=3,max=100" gorm:"size:100"`
}
