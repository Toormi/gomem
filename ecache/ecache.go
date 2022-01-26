package goecache

import (
	"github.com/orca-zhang/ecache"
	"time"
)

func NewEcache() *ecache.Cache {
	cache := ecache.NewLRUCache(16, 200, 10*time.Second).LRU2(1024)
	return cache
}
