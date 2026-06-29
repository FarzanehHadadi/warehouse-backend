package models

import (
	"encoding/json"
	"time"
)

type Activity struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	UserID      uint            `json:"user_id"`
	User        *User           `json:"user" gorm:"foreignKey:UserID"`
	Action      string          `json:"action"`
	EntityType  string          `json:"entity_type"`
	EntityID    uint            `json:"entity_id"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	Metadata    json.RawMessage `json:"metadata"`
}
