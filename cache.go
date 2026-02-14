package lrucache

type Cache[K comparable, V any] interface {
	Put(k K, v V)
	Get(k K) V
}
