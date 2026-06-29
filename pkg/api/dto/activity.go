package dto

import "time"

type SuccessActivityResponse struct {
	Data []Activity `json:"data"`
}
type Activity struct {
	ID          uint        `json:"id"`
	User        UserSummary `json:"user"`
	Action      string      `json:"action"`
	EntityType  string      `json:"entity_type"`
	EntityID    uint        `json:"entity_id"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
}
type UserSummary struct {
	ID       uint   `json:"id"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}
