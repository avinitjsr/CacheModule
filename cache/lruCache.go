package cache

import (
	"container/list"
	"errors"
)

type lruCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

// entry is a struct that holds a key-value pair.
type entry struct {
	key   string
	value interface{}
}

// NewLRUCache creates a new LRU cache with the given capacity.
func NewLRUCache(capacity int) CacheInterface {
	if capacity <= 0 {
		panic("Capacity must be greater than 0")
	}
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves an item from the cache.
func (c *lruCache) Get(key string) (interface{}, error) {
	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		return el.Value.(*entry).value, nil
	}
	return nil, errors.New("item not found in cache")
}

// Set adds an item to the cache.
func (c *lruCache) Set(key string, value interface{}) error {
	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		el.Value.(*entry).value = value
	} else {
		if c.order.Len() >= c.capacity {
			c.evict()
		}
		e := &entry{key: key, value: value}
		el := c.order.PushFront(e)
		c.items[key] = el
	}
	return nil
}

// Remove removes an item from the cache.
func (c *lruCache) Remove(key string) error {
	if el, ok := c.items[key]; ok {
		c.order.Remove(el)
		delete(c.items, key)
		return nil
	}
	return errors.New("item not found in cache")
}

// Clear clears all items from the cache.
func (c *lruCache) Clear() error {
	c.items = make(map[string]*list.Element)
	c.order.Init()
	return nil
}

// evict removes the least recently used item from the cache.
func (c *lruCache) evict() {
	el := c.order.Back()
	if el != nil {
		c.order.Remove(el)
		delete(c.items, el.Value.(*entry).key)
	}
}
