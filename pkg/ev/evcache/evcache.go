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

type Marshaler interface {
	Get(key interface{}, returnObj interface{}) (interface{}, error)
	Set(key, object interface{}, options *store.Options) error
	Delete(key interface{}) error
	Invalidate(options store.InvalidateOptions) error
	Clear() error
}

type MarshallerReturnObj func() interface{}

// Generate adapter for marshaler.Marshaler
func NewCacheMarshaller(marshaller Marshaler, returnObj MarshallerReturnObj, option *store.Options) Interface {
	return &gocacheMarshallerAdapter{
		marshaller: marshaller,
		returnObj:  returnObj,
		option:     option,
	}
}

type gocacheMarshallerAdapter struct {
	marshaller Marshaler
	returnObj  MarshallerReturnObj
	option     *store.Options
}

func (c *gocacheMarshallerAdapter) Get(key interface{}) (interface{}, error) {
	return c.marshaller.Get(key, c.returnObj())
}

func (c *gocacheMarshallerAdapter) Set(key, object interface{}) error {
	return c.marshaller.Set(key, object, c.option)
}
