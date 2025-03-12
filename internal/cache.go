package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mutex sync.RWMutex
	cache map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return item.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mutex.Lock()
		for key, val := range c.cache {
			if time.Since(val.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}
