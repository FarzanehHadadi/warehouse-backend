package repository

import (
	"errors"
	"strings"
	"time"

	"warehouse/pkg/cache"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrDuplicateKey = errors.New("Duplicate Key")
	ErrNotFound     = errors.New("item not found")
	databaseTimeout = time.Second * 30
)

func NewRepository(db *gorm.DB, c cache.Cache) *Repository {
	rc := NewRepoCache(c)

	return &Repository{
		User:       NewUserRepository(db),
		Category:   NewCategoryRepository(db, rc),
		Unit:       NewUnitRepository(db, rc),
		Department: NewDepartmentRepository(db),
		Manager:    NewManagerRepository(db),
		Product:    NewProductRepository(db),
		Store:      NewStoreRepository(db),
		Order:      NewOrderRepository(db),
		Report:     NewReportRepository(db),
		Activity:   NewActivityRepository(db),
		Dashboard:  NewDashboardRepository(db, rc),
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// Method 1: Check SQLSTATE 23505 (best for Postgres)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	// Method 2: Fallback string check
	return strings.Contains(strings.ToLower(err.Error()), "duplicate key") ||
		strings.Contains(err.Error(), "SQLSTATE 23505")
}
