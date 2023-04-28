package simplecache_test

import (
	"testing"
	"time"

	"github.com/ramitmittal/simplecache"
)

type Item struct {
	prop1 int32
	prop2 int32
}

func (c Item) Equals(ci Item) bool {
	return c.prop1 == ci.prop1 && c.prop2 == ci.prop2
}

func TestRetrieval(t *testing.T) {
	key := "itemz"
	originalItem := Item{prop1: 5, prop2: 6}

	c := simplecache.New[string, Item](time.Hour)
	c.Add(key, originalItem)

	if retrievedItem, prs := c.Get(key); !prs {
		t.Error("Expected item to be present in cache")
	} else if retrievedItem != originalItem {
		t.Error("Expected retrieved item to be equal to original item")
	}
}

func TestExpiry(t *testing.T) {
	key := "itemz"
	originalItem := Item{prop1: 5, prop2: 6}

	c := simplecache.New[string, Item](time.Microsecond)
	c.Add(key, originalItem)

	<-time.After(2 * time.Microsecond)

	if _, prs := c.Get(key); prs {
		t.Error("Expected item to not be present in cache")
	}
}

func TestEviction(t *testing.T) {
	key := "itemz"
	originalItem := Item{prop1: 5, prop2: 6}

	c := simplecache.New[string, Item](time.Microsecond)
	c.Add(key, originalItem)

	i := 1
	for i < 1024 {
		<-time.After(1 * time.Microsecond)
		i = i * 2

		if c.Len() == 0 {
			return
		}
	}

	t.Error("Expected all items from cache to be evicted")
}
