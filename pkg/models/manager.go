package models

type Manager struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" binding:"required,min=3,max=100" gorm:"not null;size:100;"`
	Phone string `json:"phone" gorm:"size:11" binding:"omitempty,min=7,max=11"`
}

type ManagerUpdate struct {
	Name  *string `json:"name" binding:"omitempty,min=3,max=100" `
	Phone *string `json:"phone" gorm:"size:11" binding:"omitempty,min=7,max=11"`
}
