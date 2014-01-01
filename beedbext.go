package beedb

import (
	"encoding/binary"
	//"hash"
	"hash/fnv"
	///"io"
	"reflect"
	"time"
	"unsafe"
)

type CacheProvider interface {
	Get(key interface{}) (interface{}, bool, error)
	Set(key interface{}, exp int, value interface{}) error
	Clear(key interface{})
	ClearAll()
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

func (cache localCache) Clear(key interface{}) {
	if _, isPresent := cache.cached[key]; isPresent {
		delete(cache.cached, key)
	}
	if _, ok := cache.exps[key]; ok {
		delete(cache.exps, key)
	}
}

func (cache localCache) ClearAll() {
	cache.cached = map[interface{}]interface{}{}
	cache.exps = map[interface{}]time.Time{}
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
	ClearAll()
	ClearObj(tableName string, key interface{})
	ClearQuery(queryStrig string, args ...interface{})
	SetCacheProvider(cacheProvider CacheProvider)
}

//default cache manager struct
type OrmCacheManager struct {
	cacheProvider CacheProvider
}

func (this OrmCacheManager) GetObj(tableName string, key interface{}) (interface{}, bool, error) {
	kHash := HashAnyList(tableName, key)
	return this.cacheProvider.Get(kHash)
}

func (this OrmCacheManager) GetQuery(queryString string, args ...interface{}) (interface{}, bool, error) {
	args = append(args, queryString)
	kHash := HashAnyList(args...)
	return this.cacheProvider.Get(kHash)
}

func (this OrmCacheManager) SetObj(tableName string, key interface{}, exp int, value interface{}) error {
	kHash := HashAnyList(tableName, key)
	return this.cacheProvider.Set(kHash, exp, value)
}

func (this OrmCacheManager) SetQuery(queryString string, exp int, value interface{}, args ...interface{}) error {
	args = append(args, queryString)
	kHash := HashAnyList(args...)
	return this.cacheProvider.Set(kHash, exp, value)
}

func (this OrmCacheManager) ClearAll() {
	this.cacheProvider.ClearAll()
}

func (this OrmCacheManager) ClearObj(tableName string, key interface{}) {
	kHash := HashAnyList(tableName, key)
	this.cacheProvider.Clear(kHash)
}

func (this OrmCacheManager) ClearQuery(queryString string, args ...interface{}) {
	args = append(args, queryString)
	kHash := HashAnyList(args...)
	this.cacheProvider.Clear(kHash)
}

func (this OrmCacheManager) SetCacheProvider(cacheProvider CacheProvider) {
	this.cacheProvider = cacheProvider
}

//default cache manager factory function
func NewOrmCacheManager(cacheProvider CacheProvider) *OrmCacheManager {
	ormCacheManager := new(OrmCacheManager)
	ormCacheManager.SetCacheProvider(cacheProvider)
	return ormCacheManager
}

var globalCacheProvider CacheProvider = newLocalCache()
var globalCacheManager CacheManager = NewOrmCacheManager(globalCacheProvider)

//hash function
func HashAnyList(args ...interface{}) uint32 {
	var hashValue uint32
	hashValue = 0
	for _, arg := range args {
		hashValue = (hashValue * 397) ^ (HashAny(arg))
	}
	return hashValue
}

var timeType reflect.Type = reflect.TypeOf(time.Now())

func HashAny(arg interface{}) uint32 {
	a := arg
	h := fnv.New32a()
	v := reflect.ValueOf(arg)
	if reflect.TypeOf(arg).Kind() == reflect.Ptr {
		a = v.Elem().Interface()
	}

	switch {
	case v.Kind() == reflect.Int:
		if unsafe.Sizeof(arg) == 32 {
			a = int32(a.(int))
		} else {
			a = int64(a.(int))
		}
	case v.Kind() == reflect.Uint:
		if unsafe.Sizeof(arg) == 32 {
			a = uint32(a.(uint))
		} else {
			a = uint64(a.(uint))
		}
	case v.Kind() <= reflect.Complex128:
	case v.Kind() == reflect.String:
		a = []byte(a.(string))
	case v.Kind() == reflect.Struct:
		if reflect.TypeOf(arg) == timeType {
			str := a.(time.Time).Format("2006-01-02 15:04:05.000 -0700")
			a = []byte(str)
		}
	default:
		return 0
		//case v.Kind() == reflect.Struct:
		//case reflect.Array:
		//case reflect.Chan:
		//case reflect.Func:
		//case reflect.Interface:
		//case reflect.Map:
		//case reflect.Ptr:
		//case reflect.Slice:
		//case reflect.UnsafePointer:
	}
	binary.Write(h, binary.LittleEndian, a)
	return h.Sum32()
}
