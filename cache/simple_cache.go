package cache

type Getter interface {
	Get(key string) interface{}
}

type SimpleCache struct {
	cache  *safeCache
	getter Getter
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}
func NewSimpleCache(getter Getter, cache Cache) *SimpleCache {
	return &SimpleCache{
		cache:  newSafeCache(cache),
		getter: getter,
	}
}

func (s *SimpleCache) Get(key string) interface{} {
	val := s.cache.get(key)
	if val != nil {
		return val
	}
	if s.getter != nil {
		val = s.getter.Get(key)
		if val == nil {
			return nil
		}
		s.cache.set(key, val)
		return val
	}
	return nil
}
func (s *SimpleCache) Set(key string, val interface{}) {
	if val == nil {
		return
	}
	s.cache.set(key, val)
}
func (s *SimpleCache) Stat() *Stat {
	return s.cache.stat()
}
