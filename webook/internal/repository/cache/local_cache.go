package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type LocalCodeCache struct {
	lock       sync.RWMutex
	localCache *cache.Cache
}

func NewLocalCodeCache(cache *cache.Cache) CodeCache {
	return &LocalCodeCache{localCache: cache, lock: sync.RWMutex{}}
}

func (c *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	c.lock.RLock()
	_, ok := c.localCache.Get(c.key(biz, phone))
	c.lock.RUnlock()
	if !ok {
		c.lock.Lock()
		c.localCache.Set(c.key(biz, phone), code, time.Minute*10)
		c.lock.Unlock()
	}
	return nil
}

func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	c.lock.Lock()
	verifyCode, ok := c.localCache.Get(c.key(biz, phone))
	c.localCache.Delete(c.key(biz, phone))
	c.lock.Unlock()
	if !ok {
		return false, nil
	}
	if verifyCode != code {
		return false, nil
	}
	return true, nil
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
