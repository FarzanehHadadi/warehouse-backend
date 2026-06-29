package repository

import (
	"time"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats() (*models.DashboardStats, error) {
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
