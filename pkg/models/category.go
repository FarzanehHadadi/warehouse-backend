package models

type Category struct {
	Name string `json:"name" gorm:"size:100;not null;uniqueIndex;"`
	ID   uint   `json:"id" gorm:"primaryKey"`
}
