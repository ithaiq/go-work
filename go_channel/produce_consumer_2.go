package main

import (
	"fmt"
	"sync"
)

var pCount int = 0
func main() {
	product := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go producter2(product, &wg)
		go Consumer2(product, &wg)
	}
	wg.Wait()
}
func producter2(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch <- pCount
		pCount++
	}
	// 生产结束了
	if pCount == 100 {
		close(ch)
	}
}
func Consumer2(ch <-chan int, wg *sync.WaitGroup) {
	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("消费者消费了产品:", v)
			} else { // ch 已经关闭了，以后都不会再生产了
				fmt.Println("所有生产的商品都已经消费完了")
				wg.Done()
				return
			}
		}
	}
}
