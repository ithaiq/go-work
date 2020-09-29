package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// 假设商品总量为2
var lCount int = 2
// 假设中奖率为 2%
var lRate int = 2
var lOut chan int = make(chan int)
// 模拟一次抽奖请求
func LGet(id int,wg *sync.WaitGroup) <-chan int{
	defer wg.Done()
	lucky := rand.Int()%100
	if lucky < lRate {
		lOut <- id
	}
	return lOut
}
// 模拟发放奖品
func LPrize(ch <- chan int) {
	for {
		select {
		case value,_:=<-ch:
			if lCount > 0 {
				lCount--
				fmt.Println("id为:",value,"的用户获奖了")
			} else {
				// 为了方便查结果这里就不打印了
			}
		}
	}
}

func main(){
	var wg sync.WaitGroup
	go LPrize(lOut)
	for i:=0;i < 100;i++ {
		wg.Add(1)
		go LGet(i,&wg)
	}
	wg.Wait()
	defer close(lOut)
}