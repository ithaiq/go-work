package main

import (
	"context"
	"fmt"
	"sync"
)

var tCount int = 0

var mutex sync.RWMutex

func main() {
	product := make(chan int)
	ctx, cancel := context.WithCancel(context.TODO())
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	for i := 0; i < 10; i++ {
		go producterT(product, ctx)
		go ConsumerT(product, ctx, cancel)
	}
	<-ctx.Done()
}
func producterT(ch chan<- int, ctx context.Context) {
	for i := 0; i < 10; i++ {
		ch <- tCount
		tCount++
	}
	// 生产结束了
	if tCount == 100 {
		close(ch)
	}
}
func ConsumerT(ch <-chan int, ctx context.Context, cancelFunc context.CancelFunc) {
	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("消费者消费了产品:", v)
			} else { // ch 已经关闭了，以后都不会再生产了
				fmt.Println("所有生产的商品都已经消费完了")
				cancelFunc()
				return
			}
		}
	}
}
