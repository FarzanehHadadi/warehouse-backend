package repository

import (
	"fmt"
	"time"
)

const (
	cachePrefixUnit        = "unit"
	cachePrefixCategory    = "category"
	cacheKeyDashboardStats = "dashboard:stats"

	cacheTTLDefault        = 5 * time.Minute
	cacheTTLDashboardStats = time.Minute
	cacheTTLActivityRecent = time.Minute
)

func cacheKeyByID(prefix string, id uint) string {
	return fmt.Sprintf("%s:id:%d", prefix, id)
}

func cacheKeyListVersion(prefix string) string {
	return fmt.Sprintf("%s:list:version", prefix)
}
