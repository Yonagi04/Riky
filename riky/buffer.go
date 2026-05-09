package riky

import (
	"strconv"
	"sync"
)

// Buffer 定义编码器需要的缓冲区接口
// 允许 encoder 包不依赖外部包，同时保持高性能
type Buffer interface {
	Bytes() []byte
	AppendByte(byte)
	AppendBytes([]byte)
	AppendString(string)
	AppendInt(int64)
	AppendFloat(float64, byte, int, int)
	Reset()
}

// buffer 是 Buffer 接口的实现
type buffer struct {
	bs []byte
}

func newBuffer() *buffer {
	return &buffer{bs: make([]byte, 0, 1024)}
}

func (b *buffer) AppendByte(c byte) {
	b.bs = append(b.bs, c)
}

func (b *buffer) AppendBytes(bs []byte) {
	b.bs = append(b.bs, bs...)
}

func (b *buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

func (b *buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

func (b *buffer) AppendFloat(f float64, fmt byte, prec, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, fmt, prec, bitSize)
}

func (b *buffer) Reset() {
	b.bs = b.bs[:0]
}

func (b *buffer) Bytes() []byte {
	return b.bs
}

// Buffer 对象池
var _bufferPool = sync.Pool{
	New: func() interface{} {
		return newBuffer()
	},
}

// GetBuffer 从池中获取 Buffer
func GetBuffer() Buffer {
	return _bufferPool.Get().(Buffer)
}

// PutBuffer 将 Buffer 归还到池中
func PutBuffer(b Buffer) {
	b.Reset()
	_bufferPool.Put(b)
}
