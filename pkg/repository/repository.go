package repository

import (
	"errors"
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
		User:     NewUserRepository(db),
		Category: NewCategoryRepository(db),
		Unit:     NewUnitRepository(db),
	}
}
func isDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
