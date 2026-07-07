package repository

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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

type cacheSettings struct {
	repo   *RepoCache
	prefix string
	ttl    time.Duration
}

type BaseRepository[T any] struct {
	db    *gorm.DB
	cache *cacheSettings
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func NewCachedBaseRepository[T any](db *gorm.DB, rc *RepoCache, prefix string, ttl time.Duration) *BaseRepository[T] {
	repo := NewBaseRepository[T](db)
	if rc != nil && rc.enabled() {
		repo.cache = &cacheSettings{repo: rc, prefix: prefix, ttl: ttl}
	}
	return repo
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

	if r.cache != nil {
		r.cache.repo.bumpListVersion(r.cache.prefix)
	}

	return nil
}

// ==================== Read ====================

func (r *BaseRepository[T]) FindByID(id uint, preloads ...string) (*T, error) {

	entity, err := r.findByIDFromDB(id, preloads...)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) findByIDFromDB(id uint, preloads ...string) (*T, error) {
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

type cachedListResult[T any] struct {
	Results []*T
	Cursor  *filter.CursorResponse
}

func (r *BaseRepository[T]) GetList(req filter.Request, config filter.FilterConfig, preloads ...string) ([]*T, *filter.CursorResponse, error) {
	if r.cache != nil {
		key := r.listCacheKey(req)
		var cached cachedListResult[T]
		if err := r.cache.repo.Get(key, &cached); err == nil {
			return cached.Results, cached.Cursor, nil
		}
	}

	results, cursor, err := r.getListFromDB(req, config, preloads...)
	if err != nil {
		return nil, nil, err
	}

	if r.cache != nil {
		key := r.listCacheKey(req)
		_ = r.cache.repo.Set(key, cachedListResult[T]{
			Results: results,
			Cursor:  cursor,
		}, r.cache.ttl)
	}

	return results, cursor, nil
}

func (r *BaseRepository[T]) getListFromDB(req filter.Request, config filter.FilterConfig, preloads ...string) ([]*T, *filter.CursorResponse, error) {
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

// listCacheRequest captures every query dimension that affects paginated list results.
type listCacheRequest struct {
	Filters []filter.Condition `json:"filters"`
	Search  string             `json:"search"`
	SortBy  string             `json:"sort_by"`
	Desc    bool               `json:"desc"`
	Cursor  string             `json:"cursor"`
	Limit   int                `json:"limit"`
}

func normalizeListCacheRequest(req filter.Request) listCacheRequest {
	return listCacheRequest{
		Filters: req.Filters,
		Search:  req.Search,
		SortBy:  req.SortBy,
		Desc:    req.Desc,
		Cursor:  req.Cursor,
		Limit:   filter.NormalizeLimit(req.Limit),
	}
}

func (r *BaseRepository[T]) listCacheKey(req filter.Request) string {
	data, _ := json.Marshal(normalizeListCacheRequest(req))
	hash := sha256.Sum256(data)
	version := r.cache.repo.listVersion(r.cache.prefix)
	return fmt.Sprintf("%s:list:%d:%x", r.cache.prefix, version, hash[:8])
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

	if r.cache != nil {
		r.cache.repo.invalidateEntity(r.cache.prefix, id)
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

	if r.cache != nil {
		r.cache.repo.invalidateEntity(r.cache.prefix, id)
	}

	return nil
}
