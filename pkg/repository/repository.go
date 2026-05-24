package repository

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrDuplicateKey = errors.New("Duplicate Key")
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Category: NewCategoryRepository(db),
	}
}
func isDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
