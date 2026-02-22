package encoder

import (
	"time"
)

// Buffer 定义编码器需要的缓冲区接口
// 允许 encoder 包不依赖 logger 包，同时保持高性能
type Buffer interface {
	Bytes() []byte
	AppendByte(byte)
	AppendBytes([]byte)
	AppendString(string)
	AppendInt(int64)
	AppendFloat(float64, byte, int, int)
}

// FieldType 表示字段类型
type FieldType int

const (
	Int64Type FieldType = iota
	StringType
	BoolType
	Float64Type
	TimeType
	DurationType
	ErrorType
)

// Field 表示一个日志字段
type Field struct {
	Key     string
	Type    FieldType
	Integer int64
	String  string
	Float   uint64
}

// Encoder 定义日志编码器接口
type Encoder interface {
	Encode(buf Buffer, msg string, level string, t time.Time, fields []Field)
	Clone() Encoder
	AddFields([]Field)
}
