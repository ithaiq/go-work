package nsqutil

const minQueueLen = 16

type Queue struct {
	buf               []*NsqdClient
	head, tail, count int
}

func NewQueue() *Queue {
	return &Queue{
		buf: make([]*NsqdClient, minQueueLen),
	}
}

func (q *Queue) Length() int {
	return q.count
}

func (q *Queue) Add(elem *NsqdClient) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem

	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

func (q *Queue) Get() *NsqdClient {
	if q.count <= 0 {
		return nil
	}
	return q.buf[q.head]
}

func (q *Queue) Remove() *NsqdClient {
	if q.count <= 0 {
		return nil
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

func (q *Queue) resize() {
	newBuf := make([]*NsqdClient, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *Queue) Find(key string) bool {
	for _, v := range q.buf {
		if v != nil && v.nsqdAddr == key {
			return true
		}
	}
	return false
}

func (q *Queue) PrintClients() string {
	var str string
	for _, v := range q.buf {
		if v != nil {
			str += v.nsqdAddr
			str += " "
		}

	}
	return str
}
