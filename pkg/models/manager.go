package models

type Manager struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required,min=3, max=100" gorm:"not null;min:3;max:100;"`
	Phone       string `json:"phone" gorm:"min:7;max:11" binding:"omitempty,min=7,max=11"`
	Departments []Department
}
