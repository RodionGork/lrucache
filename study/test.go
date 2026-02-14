package main

import (
    "fmt"
    "math/rand"
    "time"

	"github.com/rodiongork/lrucache"
)

const testSize = 20000000

func test(cache lrucache.Cache[string, float64], msg string) {
    fmt.Println(msg)
    t0 := time.Now().UnixMilli()
	r := rand.New(rand.NewSource(17))
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
	fmt.Printf("\thit: %dk, miss: %dk, time: %.1f\n", hit/1000, miss/1000, float64(dt)/1000)
}

func main() {
    for _, sz := range []int{50000, 100000, 250000, 500000} {
        suffix := fmt.Sprintf("[%dk]", sz / 1000)
        test(lrucache.NewListCache[string, float64](sz), "ListCache" + suffix)
        test(lrucache.NewGensCache[string, float64](sz), "GensCache" + suffix)
    }
}
