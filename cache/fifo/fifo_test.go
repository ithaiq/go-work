package fifo_test

import (
	"cache/fifo"
	"fmt"
	"github.com/matryer/is"
	"testing"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)
	cache := fifo.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)
	cache.Del("k1")
	is.Equal(1, cache.Len())
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)
	keys := make([]string, 0, 8)
	onRemoved := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := fifo.New(16, onRemoved)

	cache.Set("k1", 1)
	fmt.Println(cache.Len())
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}
	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}
