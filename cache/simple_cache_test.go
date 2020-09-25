package cache_test

import (
	"cache"
	"cache/fast"
	"cache/lru"
	"github.com/allegro/bigcache/v2"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleCache_Get(t *testing.T) {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}
	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("From DB find key", key)

		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})
	simpleCache := cache.NewSimpleCache(getter, lru.New(0, nil))

	is := is.New(t)
	var wg sync.WaitGroup
	for k, v := range db {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			is.Equal(simpleCache.Get(k), v)
			is.Equal(simpleCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()
	is.Equal(simpleCache.Get("unknown"), nil)
	is.Equal(simpleCache.Get("unknown"), nil)

	is.Equal(simpleCache.Stat().NGet, 10)
	is.Equal(simpleCache.Stat().NHit, 4)
}

func TestBigCache(t *testing.T) {
	bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Second))
	if err != nil {
		t.Log(err)
	}
	entry, err := bigCache.Get("my-unique-key")
	if err != nil {
		t.Log(err)
	}
	if entry == nil {
		entry = []byte("value")
		bigCache.Set("my-unique-key", entry)
	}
	t.Log(string(entry))
	entry, err = bigCache.Get("my-unique-key")
	if err != nil {
		t.Log(err)
	}
	t.Log(string(entry))
}

func BenchmarkSimpleCache(b *testing.B) {
	cache := cache.NewSimpleCache(nil, lru.New(100, nil))
	rand.Seed(time.Now().Unix())
	//require.NoError()
	b.StartTimer()
	b.Log(b.N)
	var i int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set(strconv.Itoa(rand.Intn(1000)), 1)
			atomic.AddInt32(&i, 1)
		}
	})
	b.Log(i)
	b.StopTimer()
}

func BenchmarkFastCache(b *testing.B) {
	cache := fast.NewFastCache(b.N, 100, nil)
	rand.Seed(time.Now().Unix())

	b.StartTimer()
	b.Log(b.N)
	var i int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set(strconv.Itoa(rand.Intn(1000)), 1)
			atomic.AddInt32(&i, 1)
		}
	})
	b.Log(i)
	b.StopTimer()
}
