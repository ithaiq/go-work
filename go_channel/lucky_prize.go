package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 假设商品总量为2
var count int = 2
// 假设中奖率为 2%
var rate int = 2
// 这里我们的out需要带缓冲区吗
var out chan int = make(chan int)
// 模拟一次抽奖请求
func Get(id int) <-chan int{
	// 这里能不能 close(out) 呢？
	lucky := rand.Int()%100
	if lucky < rate {
		out <- id
	}
	return out
}
// 模拟发放奖品
func Prize(ch <- chan int) {
	for {
		select {
		case value,_:=<-ch:
			if count > 0 {
				count--
				fmt.Println("id为:",value,"的用户获奖了")
			} else {
				// 为了方便查结果这里就不打印了
			}
		}
	}
}

func main(){
	go Prize(out)
	// 开启一万条协程，模拟抽奖活动
	for i:=0;i < 10000;i++ {
		go Get(i)
	}
	defer close(out)
	time.Sleep(time.Second)
}