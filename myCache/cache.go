package MyCache

import (
	"sync"

	"goTool/myCache/lru"
)

type Cache struct {
	mu  sync.RWMutex
	lru *lru.Cache

	Size int
}

func (c *Cache) Add(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = &lru.Cache{}
	}

	if c.lru.Contains(key) {
		c.lru.Remove(key)
	}
	c.lru.Add(key, value)

	if c.lru.Len() > c.Size && c.lru != nil {
		c.lru.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	vi, ok := c.lru.Get(key)
	if !ok {
		return
	}
	return vi, true
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru != nil {
		c.lru.Remove(key)
	}
}

func (c *Cache) RemoveOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru != nil {
		c.lru.RemoveOldest()
	}
}

func (c *Cache) Items() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ItemsLocked()
}

func (c *Cache) ItemsLocked() int64 {
	if c.lru == nil {
		return 0
	}
	return int64(c.lru.Len())
}
