package evcache

import (
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
)

// Cache interface
type Interface interface {
	Get(key interface{}) (interface{}, error)
	Set(key, object interface{}) error
}

// Generate adapter for cache.CacheInterface
func NewCache(cache cache.CacheInterface, option *store.Options) Interface {
	return &gocacheAdapter{
		cache:  cache,
		option: option,
	}
}

type gocacheAdapter struct {
	cache  cache.CacheInterface
	option *store.Options
}

func (c *gocacheAdapter) Get(key interface{}) (interface{}, error) {
	return c.cache.Get(key)
}

func (c *gocacheAdapter) Set(key, object interface{}) error {
	return c.cache.Set(key, object, c.option)
}
