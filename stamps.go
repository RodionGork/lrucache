package lrucache

import (
	"os"
	"sync"
)

type valueWithStamp[V any] struct {
	value V
	stamp int
}

type StampsCache[K comparable, V any] struct {
	sync.Mutex
	kv   map[K]*valueWithStamp[V]
	max  int
	fpt  int
	cur  int
	zero V
}

func NewStampsCache[K comparable, V any](size int, flushPart int) Cache[K, V] {
	return &StampsCache[K, V]{
		kv:  map[K]*valueWithStamp[V]{},
		max: size,
		fpt: flushPart,
	}
}

func (cache *StampsCache[K, V]) Put(key K, value V) {
	cache.Lock()
	defer cache.Unlock()
	cache.cur++
	if len(cache.kv) == cache.max {
		cache.flush()
	}
	if v, ok := cache.kv[key]; ok {
		v.value = value
		v.stamp = cache.cur
	} else {
		cache.kv[key] = &valueWithStamp[V]{value: value, stamp: cache.cur}
	}
}

func (cache *StampsCache[K, V]) Get(key K) V {
	cache.Lock()
	defer cache.Unlock()
	value := cache.kv[key]
	if value == nil {
		return cache.zero
	}
	cache.cur++
	value.stamp = cache.cur
	return value.value
}

func (cache *StampsCache[K, V]) flush() {
	probe := make([]int, cache.fpt)
	i := 0
	for _, v := range cache.kv {
		j := i
		for j > 0 && probe[j-1] > v.stamp {
			probe[j] = probe[j-1]
			j--
		}
		probe[j] = v.stamp
		i++
		if i == len(probe) {
			break
		}
	}
	threshold := probe[1]
	evicted := 0
	for k, v := range cache.kv {
		if v.stamp < threshold {
			delete(cache.kv, k)
			evicted++
		}
	}
	if os.Getenv("LRU_DEBUG") != "" {
		println("evicted:", evicted)
	}
}
