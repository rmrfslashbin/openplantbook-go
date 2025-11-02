package openplantbook

import (
	"testing"
	"time"
)

func TestInMemoryCache_GetSet(t *testing.T) {
	cache := NewInMemoryCache()
	defer cache.Close()

	key := "test-key"
	value := []byte("test-value")
	ttl := 1 * time.Hour

	// Get non-existent key
	_, ok := cache.Get(key)
	if ok {
		t.Error("Get() returned true for non-existent key")
	}

	// Set value
	cache.Set(key, value, ttl)

	// Get existing key
	got, ok := cache.Get(key)
	if !ok {
		t.Fatal("Get() returned false for existing key")
	}

	if string(got) != string(value) {
		t.Errorf("Get() = %q, want %q", got, value)
	}
}

func TestInMemoryCache_Expiration(t *testing.T) {
	cache := NewInMemoryCache()
	defer cache.Close()

	key := "test-key"
	value := []byte("test-value")
	ttl := 100 * time.Millisecond

	// Set value with short TTL
	cache.Set(key, value, ttl)

	// Should exist immediately
	_, ok := cache.Get(key)
	if !ok {
		t.Error("Get() returned false immediately after Set()")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	_, ok = cache.Get(key)
	if ok {
		t.Error("Get() returned true after TTL expiration")
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	cache := NewInMemoryCache()
	defer cache.Close()

	key := "test-key"
	value := []byte("test-value")
	ttl := 1 * time.Hour

	cache.Set(key, value, ttl)

	// Verify it exists
	_, ok := cache.Get(key)
	if !ok {
		t.Fatal("Get() returned false after Set()")
	}

	// Delete it
	cache.Delete(key)

	// Verify it's gone
	_, ok = cache.Get(key)
	if ok {
		t.Error("Get() returned true after Delete()")
	}
}

func TestInMemoryCache_Clear(t *testing.T) {
	cache := NewInMemoryCache()
	defer cache.Close()

	ttl := 1 * time.Hour

	// Set multiple values
	cache.Set("key1", []byte("value1"), ttl)
	cache.Set("key2", []byte("value2"), ttl)
	cache.Set("key3", []byte("value3"), ttl)

	// Verify they exist
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 not found")
	}
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 not found")
	}

	// Clear cache
	cache.Clear()

	// Verify they're all gone
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 found after Clear()")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 found after Clear()")
	}
	if _, ok := cache.Get("key3"); ok {
		t.Error("key3 found after Clear()")
	}
}

func TestInMemoryCache_Cleanup(t *testing.T) {
	cache := NewInMemoryCache()
	defer cache.Close()

	// Set items with short TTL
	cache.Set("key1", []byte("value1"), 50*time.Millisecond)
	cache.Set("key2", []byte("value2"), 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Trigger manual cleanup
	cache.removeExpired()

	// Check that expired items were removed
	cache.mu.RLock()
	count := len(cache.items)
	cache.mu.RUnlock()

	if count != 0 {
		t.Errorf("cleanup() left %d items, want 0", count)
	}
}

func TestNoOpCache(t *testing.T) {
	cache := NewNoOpCache()

	key := "test-key"
	value := []byte("test-value")
	ttl := 1 * time.Hour

	// Set should not store
	cache.Set(key, value, ttl)

	// Get should always return false
	_, ok := cache.Get(key)
	if ok {
		t.Error("NoOpCache.Get() returned true")
	}

	// Delete should not panic
	cache.Delete(key)

	// Clear should not panic
	cache.Clear()
}
