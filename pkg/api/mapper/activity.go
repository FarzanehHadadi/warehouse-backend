package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

func ToUserSummary(user *models.User) dto.UserSummary {
	return dto.UserSummary{
		ID:       user.ID,
		Phone:    user.Phone,
		Username: user.Username,
	}
}
func ToActivitySummary(activity *models.Activity) dto.Activity {
	if activity == nil {
		return dto.Activity{}
	}
	return dto.Activity{
		ID:          activity.ID,
		User:        ToUserSummary(activity.User),
		Action:      activity.Action,
		EntityType:  activity.EntityType,
		EntityID:    activity.EntityID,
		Description: activity.Description,
		CreatedAt:   activity.CreatedAt,
	}
}
func ToActivitySummaries(activities []models.Activity) []dto.Activity {
	if activities == nil {
		return []dto.Activity{}
	}
	summaries := make([]dto.Activity, len(activities))
	for i, activity := range activities {
		summaries[i] = ToActivitySummary(&activity)
	}
	return summaries
}
