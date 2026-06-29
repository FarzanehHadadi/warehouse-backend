package events

import (
	"time"
	"warehouse/pkg/logger"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Action string

const (
	Created Action = "created"
	Updated Action = "updated"
	Deleted Action = "deleted"
)

type EntityType string

const (
	Department EntityType = "department"
	Manager    EntityType = "manager"
	Unit       EntityType = "unit"
	Product    EntityType = "product"
	Store      EntityType = "store"
	Order      EntityType = "order"
	Category   EntityType = "category"
	User       EntityType = "user"
)

type Event struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	Action      Action     `json:"action"`
	UserID      uint       `json:"user_id"`
	EntityType  EntityType `json:"entity_type"`
	EntityID    uint       `json:"entity_id"`
	Description string     `json:"description"`
	Timestamp   time.Time  `json:"timestamp"`
	Metadata    any        `json:"metadata"`
}

var Bus = NewEventBus()

func Log(c *gin.Context, action Action, entityType EntityType, entityID uint, description string, payload interface{}) {
	userID := utils.GetUserIDFromContext(c)
	event := Event{
		Action:      action,
		CreatedAt:   time.Now(),
		UserID:      userID,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		Timestamp:   time.Now(),
		Metadata:    payload,
	}
	logger.Log.Info("Event logged", zap.Any("event", event))
	Bus.Publish(event)
}
