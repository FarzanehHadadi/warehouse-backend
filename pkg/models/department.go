package models

type Department struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	Name      string   `json:"name" binding:"required,min=1,max=100" gorm:"uniqueIndex; not null;size:100"`
	ManagerID *uint    `json:"manager_id" gorm:"index"`
	Manager   *Manager `json:"manager" gorm:"foreignKey:ManagerID"`
}
type DepartmentUpdate struct {
	Name      *string `json:"name" binding:"omitempty,min=1,max=100"`
	ManagerID *uint   `json:"manager_id"`
}
