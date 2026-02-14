package lrucache

import (
	"sync"
)

type GensCache[K comparable, V any] struct {
	*sync.Mutex
	fresh map[K]V
	old   map[K]V
	max   int
	zero  V
}

func NewGensCache[K comparable, V any](size int) Cache[K, V] {
	return &GensCache[K, V]{
		Mutex: &sync.Mutex{},
		fresh: map[K]V{},
		old:   map[K]V{},
		max:   size,
	}
}

func (cache *GensCache[K, V]) Put(key K, value V) {
	cache.Lock()
	defer cache.Unlock()
	sz := len(cache.fresh)
	cache.fresh[key] = value
	if sz < len(cache.fresh) {
		cache.deleteFromOldAndFlush(key)
	}
}

func (cache *GensCache[K, V]) Get(key K) V {
	cache.Lock()
	defer cache.Unlock()
	if value, ok := cache.fresh[key]; ok {
		return value
	} else if value, ok := cache.old[key]; ok {
		cache.fresh[key] = value
		cache.deleteFromOldAndFlush(key)
		return value
	}
	return cache.zero
}

func (cache *GensCache[K, V]) deleteFromOldAndFlush(key K) {
	delete(cache.old, key)
	if len(cache.fresh) == cache.max {
		cache.old = cache.fresh
		cache.fresh = make(map[K]V, cache.max)
	}
}
