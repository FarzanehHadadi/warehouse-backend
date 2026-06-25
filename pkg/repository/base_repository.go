package repository

import (
	"context"
	"errors"
	"time"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/logger"

	"gorm.io/gorm"
)

// Cursorable allows generic cursor pagination
type Cursorable interface {
	GetID() uint
	GetCreatedAt() time.Time
}

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// ==================== Query Builder ====================

func (r *BaseRepository[T]) buildQuery(ctx context.Context, req filter.Request, config filter.FilterConfig, preloads ...string) (*gorm.DB, error) {
	query := r.db.WithContext(ctx).Model(new(T))

	query, err := filter.Apply(query, req, config)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	return query, nil
}

// ==================== Create ====================

func (r *BaseRepository[T]) Create(entity *T) error {
	result := r.db.Create(entity)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return ErrDuplicateKey
		}
		logger.Log.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// ==================== Read ====================

func (r *BaseRepository[T]) FindByID(id uint, preloads ...string) (*T, error) {
	var entity T
	query := r.db

	for _, p := range preloads {
		query = query.Preload(p)
	}

	err := query.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		logger.Log.Error(err.Error())
		return nil, err
	}
	return &entity, nil
}

// ==================== List Operations ====================

func (r *BaseRepository[T]) GetList(req filter.Request, config filter.FilterConfig, preloads ...string) ([]*T, *filter.CursorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query, err := r.buildQuery(ctx, req, config, preloads...)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}

	query, err = filter.ApplyCursor(query, req, config)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}

	var results []*T
	if err := query.Find(&results).Error; err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	hasMore := len(results) > limit
	if hasMore {
		results = results[:limit]
	}

	var nextCursor string
	if len(results) > 0 {
		last := results[len(results)-1]
		if c, ok := any(last).(Cursorable); ok {
			nextCursor = filter.EncodeCursor(c.GetID(), c.GetCreatedAt())
		}
	}

	return results, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *BaseRepository[T]) GetListNoPagination(req filter.Request, config filter.FilterConfig, preloads ...string) ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query, err := r.buildQuery(ctx, req, config, preloads...)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	query = query.Order("created_at DESC, id DESC")

	var results []*T
	if err := query.Find(&results).Error; err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return results, nil
}

// ==================== Update & Delete ====================

func (r *BaseRepository[T]) Update(id uint, updates interface{}) error {
	result := r.db.Model(new(T)).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		logger.Log.Error(result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BaseRepository[T]) Delete(id uint) error {
	result := r.db.Delete(new(T), id)
	if result.Error != nil {
		logger.Log.Error(result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
