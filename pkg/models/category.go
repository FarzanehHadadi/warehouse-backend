package models

type Category struct {
	Basic
	Name string `json:"name" gorm:"size:100;not null;uniqueIndex;" binding:"required,min=3"`
}
