package _go

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"
)

func ProcessMessage(ctx context.Context, in <-chan string) {
	for {
		select {
		case s, ok := <-in:
			if !ok {
				return
			}
			//handle
			fmt.Println("s:", s)
		case <-time.After(5 * time.Minute):
		//do something
		case <-ctx.Done():
			return
		}
	}
}

//在标准库time.After的文档中有如下一段说明：
//首先等待持续时间过去，然后在返回的channel上发送当前时间。它等效于NewTimer().C。
//在计时器被触发之前，计时器不会被垃圾收集器回收。也就是说，如果没有到5min该函数就返回了，则计时器不会被GC回收，从而导致内存泄露。
//因此在使用time.After时一定要特别注意，一般来说，建议不要使用time.After，而是使用time.NewTimer
func ProcessMessage2(ctx context.Context, in <-chan string) {
	idleDuration := 5 * time.Minute
	idleDelay := time.NewTimer(idleDuration)
	//!
	defer idleDelay.Stop()
	for {
		idleDelay.Reset(idleDuration)
		select {
		case s, ok := <-in:
			if !ok {
				return
			}
			//handle
			fmt.Println("s:", s)
		case <-idleDelay.C:
		case <-ctx.Done():
			return
		}
	}
}
func ProcessMessage3(conn net.Conn) {
	var userActive = make(chan struct{})
	go func() {
		d := 1 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		// handle
		// 用户活跃
		userActive <- struct{}{}
	}
}
