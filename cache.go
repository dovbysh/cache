package cache

import (
	"sync"
	"time"
)

type Item[V any] struct {
	value   V
	expires *time.Time
}

type Cache[K comparable, V any] struct {
	items map[K]Item[V]
	mtx   sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]Item[V]),
	}
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	item := Item[V]{
		value: value,
	}
	if ttl > 0 {
		expires := time.Now().Add(ttl)
		item.expires = &expires
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.items[key] = item
}

func (c *Cache[K, V]) Delete(key K) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.items, key)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var zero V
	item, exists := c.items[key]
	if !exists {
		return zero, false
	}
	if item.expires != nil && item.expires.Before(time.Now()) {
		delete(c.items, key)
		return zero, false
	}
	return item.value, true
}

func (c *Cache[K, V]) Keys() []K {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	keys := make([]K, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache[K, V]) Len() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return len(c.items)
}
