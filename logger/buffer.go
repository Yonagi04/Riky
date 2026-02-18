package logger

import (
	"strconv"
	"sync"
)

type Buffer struct {
	bs []byte
}

func (b *Buffer) AppendByte(c byte) {
	b.bs = append(b.bs, c)
}

func (b *Buffer) AppendBytes(bs []byte) {
	b.bs = append(b.bs, bs...)
}

func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

var _bufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{bs: make([]byte, 0, 1024)}
	},
}

func getBuffer() *Buffer {
	return _bufferPool.Get().(*Buffer)
}

func putBuffer(b *Buffer) {
	b.Reset()
	_bufferPool.Put(b)
}
