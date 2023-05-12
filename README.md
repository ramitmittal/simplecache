# simplecache

A thread-safe wrapper over generic `Map[K, V]` with automatic key expiry. For times when you need to cache something in-memory without the fluff, quick!

## Installation
```
go get github.com/ramitmittal/simplecache
```

## Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/ramitmittal/simplecache"
)

func main() {
	type useless struct {
		prop1 string
		prop2 int64
	}

	// A cache with string keys and "useless" values
	// Items are evicted after 5 minutes
	uselessCache := simplecache.New[string, *useless](5 * time.Minute)

	// Add a KV pair to cache
	uselessCache.Add("foo", &useless{prop1: "eh", prop2: 4})

	// Retrieve a value from cache
	val, prs := uselessCache.Get("foo")
	if !prs {
		fmt.Println("Item not present in cache.")
	} else {
		fmt.Printf("Found: %v\n", val)
	}
}
```
