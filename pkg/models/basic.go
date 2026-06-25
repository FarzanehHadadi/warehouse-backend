package models

import "time"

type Basic struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;default:current_timestamp;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;default:current_timestamp;autoCreateTime"`
}

func (b Basic) GetID() uint {
	return b.ID
}

func (b Basic) GetCreatedAt() time.Time {
	return b.CreatedAt
}
