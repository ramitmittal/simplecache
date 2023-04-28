package simplecache

import (
	"sync"
	"time"
)

type valueWrapper[S any] struct {
	addedOn time.Time
	value   S
}

type Cache[T comparable, S any] struct {
	cacheDuration time.Duration
	mu            sync.RWMutex
	items         map[T]valueWrapper[S]
}

func New[T comparable, S any](cacheDuration time.Duration) *Cache[T, S] {
	c := Cache[T, S]{
		cacheDuration: cacheDuration,
		items:         make(map[T]valueWrapper[S]),
	}
	go c.runMaintenanceTasks()
	return &c
}

func (c *Cache[T, S]) runMaintenanceTasks() {
	for now := range time.NewTicker(c.cacheDuration).C {
		c.evict(now)
	}
}

// Delete expired items from cache.
func (c *Cache[T, S]) evict(now time.Time) {
	c.mu.Lock()
	for k, v := range c.items {
		if v.addedOn.Add(c.cacheDuration).Before(now) {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Add an item to cache.
func (c *Cache[T, S]) Add(k T, v S) {
	c.mu.Lock()
	c.items[k] = valueWrapper[S]{
		addedOn: time.Now(),
		value:   v,
	}
	c.mu.Unlock()
}

// Get an item from cache. Returns the item and a boolean indicating whether the key was found.
// Returns the zero-value of the item if the key is not found.
func (c *Cache[T, S]) Get(k T) (S, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if vw, prs := c.items[k]; !prs {
		x := new(S)
		return *x, false
	} else if time.Now().After(vw.addedOn.Add(c.cacheDuration)) {
		x := new(S)
		return *x, false
	} else {
		return vw.value, true
	}
}

func (c *Cache[T, S]) Len() int {
	return len(c.items)
}
