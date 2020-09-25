package fast

type fastCache struct {
	shards    []*cacheShard
	shardMask uint64
	hash      fnv64a
}

func NewFastCache(maxEntries, shardsNum int, onEvicted func(key string, value interface{})) *fastCache {
	fastCache := &fastCache{
		hash:      newDefaultHasher(),
		shards:    make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
	}
	for i := 0; i < shardsNum; i++ {
		fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
	}

	return fastCache
}

func (c *fastCache) getShard(key string) *cacheShard {
	hashKey := c.hash.Sum64(key)
	return c.shards[hashKey%c.shardMask]
}

func (c *fastCache) Set(key string, value interface{}) {
	c.getShard(key).Set(key, value)
}

func (c *fastCache) Get(key string) interface{} {
	return c.getShard(key).Get(key)
}

func (c *fastCache) Del(key string) {
	c.getShard(key).Del(key)
}
func (c *fastCache) Len() int {
	length := 0
	for _, shard := range c.shards {
		length += shard.Len()
	}
	return length
}

func (c *fastCache) DelOldest() {
	panic("no implements")
}
