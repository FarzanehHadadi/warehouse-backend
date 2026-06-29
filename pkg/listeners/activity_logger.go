package listeners

import (
	"warehouse/pkg/events"
	"warehouse/pkg/logger"
	"warehouse/pkg/repository"

	"go.uber.org/zap"
)

// ActivityLogger listens to events and saves them to database
type ActivityLogger struct {
	ActivityRepo repository.ActivityRepository
}

func NewActivityLogger(repo repository.ActivityRepository) *ActivityLogger {
	return &ActivityLogger{
		ActivityRepo: repo,
	}
}

// Handle is called whenever an event is published
func (l *ActivityLogger) Handle(event events.Event) {
	err := l.ActivityRepo.Log(
		event.UserID,
		string(event.Action),
		event.EntityType,
		event.EntityID,
		event.Description,
		event.Metadata,
	)

	if err != nil {
		logger.Log.Error("failed to log activity", zap.Error(err), zap.Any("event", event))
	}
}
