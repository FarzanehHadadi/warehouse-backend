package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrDuplicateKey = errors.New("Duplicate Key")
	ErrNotFound     = errors.New("item not found")
	databaseTimeout = time.Second * 30
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       NewUserRepository(db),
		Category:   NewCategoryRepository(db),
		Unit:       NewUnitRepository(db),
		Department: NewDepartmentRepository(db),
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
