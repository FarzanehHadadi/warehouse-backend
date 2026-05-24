package dto

type Category struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
