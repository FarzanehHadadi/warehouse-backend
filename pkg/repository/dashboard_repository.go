package repository

import (
	"time"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db    *gorm.DB
	cache *RepoCache
}

func NewDashboardRepository(db *gorm.DB, rc *RepoCache) DashboardRepository {
	return &dashboardRepository{db: db, cache: rc}
}

func (r *dashboardRepository) GetStats() (*models.DashboardStats, error) {
	if r.cache != nil && r.cache.enabled() {
		var stats models.DashboardStats
		if err := r.cache.Get(cacheKeyDashboardStats, &stats); err == nil {
			return &stats, nil
		}
	}

	stats, err := r.loadStatsFromDB()
	if err != nil {
		return nil, err
	}

	if r.cache != nil && r.cache.enabled() {
		_ = r.cache.Set(cacheKeyDashboardStats, stats, cacheTTLDashboardStats)
	}

	return stats, nil
}

func (r *dashboardRepository) loadStatsFromDB() (*models.DashboardStats, error) {
	var stats models.DashboardStats
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Run multiple counts in parallel using GORM (very efficient)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Total Products
		if err := tx.Model(&models.Product{}).Count(&stats.TotalProducts).Error; err != nil {
			return err
		}

		// Orders Today
		if err := tx.Model(&models.Order{}).
			Where("created_at >= ?", todayStart).
			Count(&stats.TodayOrders).Error; err != nil {
			return err
		}

		// Active Stores
		if err := tx.Model(&models.Store{}).
			Where("is_active = ?", true).
			Count(&stats.ActiveStores).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).
			Where("inventory_count <= warning_threshold").
			Count(&stats.StoresThreshold).Error; err != nil {
			return err
		}

		return nil
	})

	return &stats, err
}
