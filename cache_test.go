package cache

import (
	"sync"
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	t.Run("basic set without TTL", func(t *testing.T) {
		c := New[string, int]()
		c.Set("key", 123, 0)

		if len(c.items) != 1 {
			t.Errorf("expected 1 item, got %d", len(c.items))
		}
		if item, exists := c.items["key"]; !exists {
			t.Error("key not found in cache")
		} else if item.value != 123 {
			t.Errorf("expected value 123, got %d", item.value)
		} else if item.expires != nil {
			t.Error("expected no expiration, but got one")
		}
	})

	t.Run("set with TTL", func(t *testing.T) {
		c := New[string, int]()
		now := time.Now()
		ttl := time.Minute
		c.Set("key", 123, ttl)

		item, exists := c.items["key"]
		if !exists {
			t.Fatal("key not found in cache")
		}
		if item.expires == nil {
			t.Error("expected expiration time, got nil")
		}
		if item.expires.Sub(now) < ttl-time.Second {
			t.Error("expiration time set incorrectly")
		}
	})

	t.Run("overwrite existing value", func(t *testing.T) {
		c := New[string, int]()
		c.Set("key", 123, 0)
		c.Set("key", 456, 0)

		if len(c.items) != 1 {
			t.Errorf("expected 1 item, got %d", len(c.items))
		}
		if item, exists := c.items["key"]; !exists {
			t.Error("key not found in cache")
		} else if item.value != 456 {
			t.Errorf("expected value 456, got %d", item.value)
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		c := New[int, int]()
		var wg sync.WaitGroup
		numGoroutines := 100

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(val int) {
				defer wg.Done()
				c.Set(val, val, 0)
			}(i)
		}
		wg.Wait()

		if len(c.items) != numGoroutines {
			t.Errorf("expected %d items, got %d", numGoroutines, len(c.items))
		}
		for i := 0; i < numGoroutines; i++ {
			if item, exists := c.items[i]; !exists {
				t.Errorf("key %d not found in cache", i)
			} else if item.value != i {
				t.Errorf("expected value %d, got %d", i, item.value)
			}
		}
	})
}
func TestCache_Get(t *testing.T) {
	t.Run("get existing item", func(t *testing.T) {
		c := New[string, int]()
		c.Set("key", 123, 0)

		val, exists := c.Get("key")
		if !exists {
			t.Error("expected item to exist")
		}
		if val != 123 {
			t.Errorf("expected value 123, got %d", val)
		}
	})

	t.Run("get non-existing item", func(t *testing.T) {
		c := New[string, int]()
		_, exists := c.Get("key")
		if exists {
			t.Error("expected item to not exist")
		}
	})

	t.Run("get expired item", func(t *testing.T) {
		c := New[string, int]()
		c.Set("key", 123, time.Millisecond)
		time.Sleep(2 * time.Millisecond)

		_, exists := c.Get("key")
		if exists {
			t.Error("expected expired item to not exist")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		c := New[int, int]()
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			c.Set(i, i, 0)
		}

		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(val int) {
				defer wg.Done()
				got, exists := c.Get(val)
				if !exists {
					t.Errorf("key %d should exist", val)
				}
				if got != val {
					t.Errorf("expected %d, got %d", val, got)
				}
			}(i)
		}
		wg.Wait()
	})
}
