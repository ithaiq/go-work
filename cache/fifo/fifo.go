package fifo

import (
	"cache"
	"container/list"
)

type fifo struct {
	maxBytes  int
	onRemoved func(key string, value interface{})
	useBytes  int
	ll        *list.List
	cache     map[string]*list.Element
}

func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.useBytes = f.useBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}
	en := &entry{key, value}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.useBytes += en.Len()
	if f.maxBytes > 0 && f.useBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}
	return nil
}

func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.useBytes -= en.Len()
	delete(f.cache, en.key)
	if f.onRemoved != nil {
		f.onRemoved(en.key, en.value)
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func New(maxBytes int, onRemoved func(key string, value interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onRemoved: onRemoved,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
