package riky

import (
	"sync/atomic"
)

// AtomicLevel 支持原子读写的日志级别
type AtomicLevel struct {
	l atomic.Int32
}

// NewAtomicLevel 创建默认级别为 InfoLevel 的 AtomicLevel
func NewAtomicLevel() *AtomicLevel {
	return &AtomicLevel{l: atomic.Int32{}}
}

// NewAtomicLevelAt 创建指定级别的 AtomicLevel
func NewAtomicLevelAt(l Level) *AtomicLevel {
	a := &AtomicLevel{l: atomic.Int32{}}
	a.l.Store(int32(l))
	return a
}

// Level 原子读取当前级别
func (a *AtomicLevel) Level() Level {
	return Level(a.l.Load())
}

// SetLevel 原子设置新级别
func (a *AtomicLevel) SetLevel(l Level) {
	a.l.Store(int32(l))
}

// Enabled 实现 LevelEnabler 接口
// 当 target >= 当前级别时返回 true
func (a *AtomicLevel) Enabled(target Level) bool {
	return target >= a.Level()
}

// String 返回当前级别的字符串表示
func (a *AtomicLevel) String() string {
	return a.Level().String()
}
