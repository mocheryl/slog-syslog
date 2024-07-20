package slogsyslog

import "sync"

// bufPool is a pool for byte slices used to create messages before being sent
// to the syslog writer.
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

// allocBuf returns a buffered byte slice.
func allocBuf() *[]byte { return bufPool.Get().(*[]byte) }

// freeBuf returns smaller byte slice back to the pool.
func freeBuf(b *[]byte) {
	if cap(*b) > maxBufferSize {
		return
	}

	*b = (*b)[:0]
	bufPool.Put(b)
}
