package repository

import (
	"context"
	"time"
	"warehouse/pkg/cache"
)

// RepoCache wraps the cache client for repository use.
type RepoCache struct {
	cache cache.Cache
}

func NewRepoCache(c cache.Cache) *RepoCache {
	return &RepoCache{cache: c}
}

func (rc *RepoCache) enabled() bool {
	return rc != nil && rc.cache != nil
}

func (rc *RepoCache) Get(key string, dest any) error {
	if !rc.enabled() {
		return cache.ErrCacheMiss
	}
	return rc.cache.Get(context.Background(), key, dest)
}

func (rc *RepoCache) Set(key string, value any, ttl time.Duration) error {
	if !rc.enabled() {
		return nil
	}
	return rc.cache.Set(context.Background(), key, value, ttl)
}

func (rc *RepoCache) Delete(key string) {
	if !rc.enabled() {
		return
	}
	_ = rc.cache.Delete(context.Background(), key)
}

func (rc *RepoCache) listVersion(prefix string) int64 {
	var version int64
	if err := rc.Get(cacheKeyListVersion(prefix), &version); err != nil {
		return 0
	}
	return version
}

func (rc *RepoCache) bumpListVersion(prefix string) {
	if !rc.enabled() {
		return
	}
	_, _ = rc.cache.Incr(context.Background(), cacheKeyListVersion(prefix))
}

func (rc *RepoCache) invalidateEntity(prefix string, id uint) {
	rc.Delete(cacheKeyByID(prefix, id))
	rc.bumpListVersion(prefix)
}
