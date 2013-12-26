package beedb

import (
	"time"
)

type CacheProvider interface {
	Get(key string) (interface{}, bool, error)
	Set(key string, exp int, value interface{}) error
}

//default cache implement class
type localCache struct {
	cached map[string]interface{}
	exps   map[string]time.Time
}

func (cache localCache) Get(key string) (value interface{}, isPresent bool, err error) {
	value, isPresent = cache.cached[key]
	err = nil
	if exp, ok := cache.exps[key]; ok {
		if time.Now().After(exp) {
			delete(cache.exps, key)
		}

	}
	return
}

func (cache localCache) Set(key string, exp int, value interface{}) error {
	cache.cached[key] = value
	if exp > 0 {
		t1 := time.Now()

		t := t1.Add(time.Second * time.Duration(exp))
		cache.exps[key] = t
	}
	return nil
}

//default cacheProvider factory function
func newLocalCache() *localCache {
	cache := new(localCache)
	cache.cached = map[string]interface{}{}
	cache.exps = map[string]time.Time{}
	return cache
}

var cacheProvider CacheProvider = newLocalCache()
