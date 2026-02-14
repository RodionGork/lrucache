package lrucache

import (
	"sync"
)

type ListElem[K any, V any] struct {
	key   K
	value V
	next  *ListElem[K, V]
	prev  *ListElem[K, V]
}

type ListCache[K comparable, V any] struct {
	*sync.Mutex
	head *ListElem[K, V]
	tail *ListElem[K, V]
	kv   map[K]*ListElem[K, V]
	max  int
	zero V
}

func NewListCache[K comparable, V any](size int) Cache[K, V] {
	cache := &ListCache[K, V]{
		Mutex: &sync.Mutex{},
		max:   size,
		kv:    map[K]*ListElem[K, V]{},
		head:  &ListElem[K, V]{},
		tail:  &ListElem[K, V]{},
	}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	return cache
}

func (cache *ListCache[K, V]) Put(key K, value V) {
	cache.Lock()
	defer cache.Unlock()
	elem, ok := cache.kv[key]
	if !ok {
		elem = &ListElem[K, V]{key: key, value: value}
		if len(cache.kv) == cache.max {
			delete(cache.kv, cache.tail.prev.key)
			detach(cache.tail.prev)
		}
		cache.kv[key] = elem
	} else {
		cache.kv[key].value = value
	}
	cache.promote(elem)
}

func (cache *ListCache[K, V]) Get(key K) V {
	cache.Lock()
	defer cache.Unlock()
	elem, ok := cache.kv[key]
	if !ok {
		return cache.zero
	}
	cache.promote(elem)
	return elem.value
}

func (cache *ListCache[K, V]) promote(elem *ListElem[K, V]) {
	if elem.prev != nil {
		detach(elem)
	}
	elem.next = cache.head.next
	elem.prev = cache.head
	elem.next.prev = elem
	elem.prev.next = elem
}

func detach[K any, V any](elem *ListElem[K, V]) {
	prev := elem.prev
	next := elem.next
	prev.next = next
	next.prev = prev
}
