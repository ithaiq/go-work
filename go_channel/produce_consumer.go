package main

import (
	"fmt"
)

func main() {
	product := make(chan int)
	wait := make(chan struct{})
	go producter(product)
	go consumer(product, wait)
	<-wait
}
func producter(ch chan<- int) {
	for i := 0; i < 50; i++ {
		fmt.Println("生产了产品", i)
		ch <- i
	}
	close(ch)
}
func consumer(ch <-chan int, wait chan<- struct{}) {
	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("消费者消费了产品:",v)
			} else { // ch 已经关闭了，以后都不会再生产了
				fmt.Println("所有生产的商品都已经消费完了")
				wait <- struct{}{}
				return
			}
		}
	}
}
