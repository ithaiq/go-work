package lfu_test

import (
	"cache/lfu"
	"github.com/matryer/is"
	"testing"
)

func TestSet(t *testing.T) {
	is := is.New(t)
	cache := lfu.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	cache.Len()
	is.Equal(v, 1)

	cache.Del("k1")
	cache.Len()
	is.Equal(0, cache.Len())
}
