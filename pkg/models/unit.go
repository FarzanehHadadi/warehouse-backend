package models

type Unit struct {
	Name string `json:"name" gorm:"size:100;not null;uniqueIndex;" binding:"required,min=1"`
	Id   uint   `json:"id" gorm:"primaryKey"`
}
