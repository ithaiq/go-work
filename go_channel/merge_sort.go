package main

import (
	"fmt"
	"sort"
)

// 我们先去接收一个序列 a... 的意思就是我要接收的参数不确定有多长，但一定是 int 型
// 操作方法和切片相同
// 思考一下在哪里进行内部排序
func GetNumers(a ...int) <-chan int {
	// 是否需要关闭 out
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)
	// 是否需要关闭 out
	go func() {
		var a []int
		for v := range in {
			a = append(a, v)
		}
		sort.Ints(a)
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok1 || (ok2 && v2 <= v1) {
				out <- v2
				v2, ok2 = <-in2
			} else {
				out <- v1
				v1, ok1 = <-in1
			}
		}
		close(out)
	}()
	return out
}

func main() {
	p1 := InMemSort(GetNumers(4, 2, 7, 5, 1, 9))
	p2 := InMemSort(GetNumers(7, 3, 6, 0, 8))
	out := Merge(p1, p2)

	for v := range out {
		fmt.Printf("%d ", v)
	}
	fmt.Printf("\n")
}
