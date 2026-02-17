package lrucache

import (
	"testing"
)

func fillCache(cache Cache[string, string]) {
	tests := [][]string{{"Abrikos", "ABR"}, {"Bananas", "BAN"}, {"Cucumba", "CUC"}, {"Abrikos"},
		{"Durian", "DUR"}, {"Eggplant", "EGG"}, {"Durian", "DUN"}, {"Fig", "FIG"}, {"Bananas"}}
	for _, v := range tests {
		if len(v) == 2 {
			cache.Put(v[0], v[1])
		} else {
			cache.Get(v[0])
		}
	}
}

func TestList(t *testing.T) {
	cache := NewListCache[string, string](4)
	fillCache(cache)
	res := ""
	lcache := cache.(*ListCache[string, string])
	for elem := lcache.head.next; elem != lcache.tail; elem = elem.next {
		res = res + elem.value + " "
	}
	if res != "FIG DUN EGG ABR " {
		t.Error("wrong list at the end: " + res)
	}
}

func TestGens(t *testing.T) {
	cache := NewGensCache[string, string](3)
	fillCache(cache)
	gcache := cache.(*GensCache[string, string])
	if len(gcache.fresh) != 2 || gcache.fresh["Fig"] != "FIG" || gcache.fresh["Durian"] != "DUN" {
		t.Errorf("wrong fresh at the end: %v", gcache.fresh)
	}
	if len(gcache.old) != 2 || gcache.old["Abrikos"] != "ABR" || gcache.old["Eggplant"] != "EGG" {
		t.Errorf("wrong old at the end: %v", gcache.fresh)
	}
}

func TestStamps(t *testing.T) {
	cache := NewStampsCache[string, string](4, 3)
	fillCache(cache)
	scache := cache.(*StampsCache[string, string])
	t.Logf("map content: %v", scache.kv)
	if len(scache.kv) > 4 {
		t.Errorf("wrong size at the end: %d", len(scache.kv))
	}
}
