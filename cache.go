package openplantbook

import (
	"sync"
	"time"
)

// Cache is the interface for caching API responses
type Cache interface {
	// Get retrieves a value from the cache
	Get(key string) ([]byte, bool)

	// Set stores a value in the cache with a TTL
	Set(key string, value []byte, ttl time.Duration)

	// Delete removes a value from the cache
	Delete(key string)

	// Clear removes all values from the cache
	Clear()
}

// InMemoryCache implements Cache using an in-memory map
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
	stop  chan struct{}
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache with background cleanup
func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		items: make(map[string]*cacheItem),
		stop:  make(chan struct{}),
	}

	// Start background cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set stores a value in the cache with a TTL
func (c *InMemoryCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
}

// Close stops the background cleanup goroutine
func (c *InMemoryCache) Close() {
	close(c.stop)
}

// cleanup removes expired items periodically
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stop:
			return
		}
	}
}

// removeExpired removes all expired items from the cache
func (c *InMemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// NoOpCache is a cache that does nothing (useful for disabling caching)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns false
func (c *NoOpCache) Get(key string) ([]byte, bool) {
	return nil, false
}

// Set does nothing
func (c *NoOpCache) Set(key string, value []byte, ttl time.Duration) {
	// No-op
}

// Delete does nothing
func (c *NoOpCache) Delete(key string) {
	// No-op
}

// Clear does nothing
func (c *NoOpCache) Clear() {
	// No-op
}
