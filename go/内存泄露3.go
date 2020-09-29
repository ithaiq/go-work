package _go

/*func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}*/
//直到 MessageChannel 这个 channel 关闭，for-range 循环才会结束，
//因此需要有地方调用close（u.MessageChannel）。这种情况的另一种情形是：
//虽然没有for-range循环，但给channel发送数据的一方已经不再发送数据了，而接收的一方还在等待，
//并且这个等待会无限期持续下去，唯一能取消它等待的方法就是关闭这个channel。
