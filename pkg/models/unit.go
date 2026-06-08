package models

type Unit struct {
	Basic
	Name string `json:"name" gorm:"size:100;not null;uniqueIndex;" binding:"required,min=1"`
}
