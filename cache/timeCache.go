package cache

import (
	"sync"
	"time"
)

type TimeCache struct {
	mu       sync.Mutex
	cache    map[int]*cacheItem
	duration time.Duration
}
type cacheItem struct {
	value      int
	expiration int64
}

func NewTimeCache(duration time.Duration) *TimeCache {
	return &TimeCache{
		cache:    make(map[int]*cacheItem),
		duration: duration,
	}
}

// Get retrieves a value from the cache
func (c *TimeCache) Get(key int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.cache[key]; found {
		if time.Now().UnixNano() < item.expiration {
			return item.value, true
		}
		delete(c.cache, key)
	}
	return 0, false
}

// Put inserts a value into the cache
func (c *TimeCache) Put(key, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(c.duration).UnixNano(),
	}
}
