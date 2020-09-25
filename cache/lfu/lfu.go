package lfu

import (
	"cache"
	"container/heap"
	"fmt"
)

type lfu struct {
	maxBytes  int
	onRemoved func(key string, value interface{})
	usedBytes int

	queue *queue
	cache map[string]*entry
}

func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
	}
	en := &entry{key: key, value: value}
	heap.Push(l.queue, en)
	l.cache[key] = en
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}
	return nil
}

func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) Len() int {
	fmt.Println("use:", l.usedBytes)
	fmt.Println("max:", l.maxBytes)
	return l.queue.Len()
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}
	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()
	if l.onRemoved != nil {
		l.onRemoved(en.key, en.value)
	}
}

func New(maxBytes int, onRemoved func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onRemoved: onRemoved,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}
