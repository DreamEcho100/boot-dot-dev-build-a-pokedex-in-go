package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		mu:       sync.Mutex{},
	}

	go cache.reapLoop()

	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cacheEntry, ok := cache.entries[key]; ok {
		return cacheEntry.val, ok
	}

	return nil, false
}

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	defer ticker.Stop()

	for range ticker.C {
		cache.mu.Lock()
		for key, cacheEntry := range cache.entries {
			if time.Since(cacheEntry.createdAt) > cache.interval {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}
