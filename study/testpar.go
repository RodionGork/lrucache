package main

import (
    "fmt"
    "math/rand"
    "time"
    "sync"

	"github.com/rodiongork/lrucache"
)

const testSize = 20000000

var wg sync.WaitGroup

func test(cache lrucache.Cache[string, float64], seed int64, msg string) {
    t0 := time.Now().UnixMilli()
	r := rand.New(rand.NewSource(seed))
	hit := 0
	miss := 0
	for i := 0; i < testSize; i++ {
		rnd := r.Float64()
		key := fmt.Sprintf("%06d", int(rnd*rnd*rnd*1000000))
		if cache.Get(key) == 0 {
			cache.Put(key, rnd)
			miss++
		} else {
			hit++
		}
	}
	dt := time.Now().UnixMilli() - t0
	fmt.Printf("%s hit: %dk, miss: %dk, time: %.1f\n", msg, hit/1000, miss/1000, float64(dt)/1000)
    wg.Done()
}

func main() {
    cache := lrucache.NewGensCache[string, float64](100000)
    wg.Add(2)
    go test(cache, 23, "GensCache 100k-23")
    go test(cache, 29, "GensCache 100k-29")
    fmt.Println("Started two tests in parallel...")
    wg.Wait()
}
