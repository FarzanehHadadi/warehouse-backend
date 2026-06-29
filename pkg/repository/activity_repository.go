package repository

import (
	"encoding/json"
	"time"

	"warehouse/pkg/logger"
	"warehouse/pkg/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// activityRepository implements ActivityRepository
type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Log(userID uint, action, entityType string, entityID uint, description string, payload interface{}) error {
	var metadata json.RawMessage
	if payload != nil {
		if b, err := json.Marshal(payload); err == nil {
			metadata = json.RawMessage(b)
		}
	}

	activity := models.Activity{
		UserID:      userID,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
	}
	logger.Log.Info("Activity logged", zap.Any("activity", activity))
	return r.db.Create(&activity).Error
}

func (r *activityRepository) GetRecent(limit int) ([]models.Activity, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var activities []models.Activity

	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, phone, first_name, last_name")
	}).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error

	return activities, err
}
