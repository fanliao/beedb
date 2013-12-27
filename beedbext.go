package beedb

import (
	"time"
)

type CacheProvider interface {
	Get(key interface{}) (interface{}, bool, error)
	Set(key interface{}, exp int, value interface{}) error
}

//default cache implement class
type localCache struct {
	cached map[interface{}]interface{}
	exps   map[interface{}]time.Time
}

func (cache localCache) Get(key interface{}) (value interface{}, isPresent bool, err error) {
	value, isPresent = cache.cached[key]
	err = nil
	if exp, ok := cache.exps[key]; ok {
		if time.Now().After(exp) {
			delete(cache.exps, key)
		}

	}
	return
}

func (cache localCache) Set(key interface{}, exp int, value interface{}) error {
	cache.cached[key] = value
	if exp > 0 {
		t1 := time.Now()

		t := t1.Add(time.Second * time.Duration(exp))
		cache.exps[key] = t
	}
	return nil
}

//default cacheProvider factory function
func newLocalCache() CacheProvider {
	cache := new(localCache)
	cache.cached = map[interface{}]interface{}{}
	cache.exps = map[interface{}]time.Time{}
	return cache
}

//cache manager
type CacheManager interface {
	GetObj(tableName string, key interface{}) (interface{}, bool, error)
	GetQuery(queryStrig string, args ...interface{}) (interface{}, bool, error)
	SetObj(tableName string, key interface{}, exp int, value interface{}) error
	SetQuery(queryStrig string, exp int, value interface{}, args ...interface{}) error
	ClearAll() error
	ClearObj(tableName string, key interface{}) error
	ClearQuery(queryStrig string, args ...interface{}) error
	SetCacheProvider(cacheProvider CacheProvider)
}

type OrmCacheManager struct {
	cacheProvider CacheProvider
}

func (this OrmCacheManager) GetObj(tableName string, key interface{}) (interface{}, bool, error) {
	return nil, true, nil
}

func (this OrmCacheManager) GetQuery(queryStrig string, args ...interface{}) (interface{}, bool, error) {
	return nil, true, nil
}

func (this OrmCacheManager) SetObj(tableName string, key interface{}, exp int, value interface{}) error {
	return nil
}

func (this OrmCacheManager) SetQuery(queryStrig string, exp int, value interface{}, args ...interface{}) error {
	return nil
}

func (this OrmCacheManager) ClearAll() error {
	return nil
}

func (this OrmCacheManager) ClearObj(tableName string, key interface{}) error {
	return nil
}

func (this OrmCacheManager) ClearQuery(queryStrig string, args ...interface{}) error {
	return nil
}

func (this OrmCacheManager) SetCacheProvider(cacheProvider CacheProvider) {
	this.cacheProvider = cacheProvider
}

func NewOrmCacheManager(cacheProvider CacheProvider) *OrmCacheManager {
	ormCacheManager := new(OrmCacheManager)
	ormCacheManager.SetCacheProvider(cacheProvider)
	return ormCacheManager
}

var globalCacheProvider CacheProvider = newLocalCache()
var globalCacheManager CacheManager = NewOrmCacheManager(globalCacheProvider)

func hash(args ...interface{}) int32 {

}
