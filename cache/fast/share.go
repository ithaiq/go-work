package fast

import (
	"cache"
	"container/list"
	"sync"
)

type cacheShard struct {
	locker     sync.RWMutex
	maxEntries int
	onRemoved  func(key string, value interface{})
	ll         *list.List
	cache      map[string]*list.Element
}

func (l *cacheShard) Set(key string, value interface{}) {
	l.locker.Lock()
	defer l.locker.Unlock()
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}
	en := &entry{key: key, value: value}
	e := l.ll.PushBack(en)
	l.cache[key] = e
	if l.maxEntries > 0 && l.ll.Len() > l.maxEntries {
		l.removeElement(l.ll.Front())
	}
}

func (l *cacheShard) Get(key string) interface{} {
	l.locker.RLock()
	defer l.locker.RUnlock()
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (l *cacheShard) Del(key string) {
	l.locker.Lock()
	defer l.locker.Unlock()

	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *cacheShard) DelOldest() {
	l.locker.Lock()
	defer l.locker.Unlock()

	l.removeElement(l.ll.Front())
}

func (l *cacheShard) Len() int {
	l.locker.RLock()
	defer l.locker.RUnlock()

	return l.ll.Len()
}

func (l *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	l.ll.Remove(e)
	en := e.Value.(*entry)
	delete(l.cache, en.key)
	if l.onRemoved != nil {
		l.onRemoved(en.key, en.value)
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func newCacheShard(maxEntries int, onRemoved func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onRemoved:  onRemoved,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}
