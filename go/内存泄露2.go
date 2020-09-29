package _go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//如果超时没有处理完，则 ctx.Done 会执行，接新旧函数返回。
//新开启的 goroutine 会因为channel中的另一端没有按时接收goroutine而一直阻塞，进而导致goroutine泄露
//这种因为发送到channel阻塞而导致goroutine泄露的最简单的解决方法是将channel改为有缓冲的channel，并保证容量充足。
//比如上面的例子，只需将ch改为：ch：=make（chan string，1）即可s
func process(term string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	ch := make(chan string)
	go func() {
		ch <- search(term)
	}()
	select {
	case <-ctx.Done():
		return errors.New("search canceled")
	case result := <-ch:
		fmt.Println("Received:", result)
		return nil
	}
}
func search(term string) string {
	time.Sleep(200 * time.Millisecond)
	return "some value"
}
