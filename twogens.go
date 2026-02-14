package lrucache

import (
	"sync"
	"sync/atomic"
)

type GensCache[K comparable, V any] struct {
	*sync.Mutex
	fresh *sync.Map
	old   *sync.Map
	sz    atomic.Int32
	max   int32
	zero  V
}

func NewGensCache[K comparable, V any](size int) Cache[K, V] {
	return &GensCache[K, V]{
		Mutex: &sync.Mutex{},
		fresh: &sync.Map{},
		old:   &sync.Map{},
		sz:    atomic.Int32{},
		max:   int32(size),
	}
}

func (cache *GensCache[K, V]) Put(key K, value V) {
	cache.Lock()
	defer cache.Unlock()
	_, replaced := cache.fresh.Swap(key, value)
	if !replaced {
		cache.sz.Add(1)
		cache.deleteFromOldAndFlush(key)
	}
}

func (cache *GensCache[K, V]) Get(key K) V {
	if value, ok := cache.fresh.Load(key); ok {
		return value.(V)
	}
	cache.Lock()
	defer cache.Unlock()
	if value, ok := cache.old.Load(key); ok {
		cache.fresh.Store(key, value)
		cache.deleteFromOldAndFlush(key)
		return value.(V)
	}
	return cache.zero
}

func (cache *GensCache[K, V]) deleteFromOldAndFlush(key K) {
	cache.old.Delete(key)
	if cache.sz.Load() >= cache.max {
		cache.old = cache.fresh
		cache.fresh = &sync.Map{}
		cache.sz.Store(0)
	}
}
