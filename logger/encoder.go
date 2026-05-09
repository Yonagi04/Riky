package logger

import (
	"time"
)

// Encoder 定义日志编码器接口
type Encoder interface {
	Encode(buf Buffer, msg string, level string, t time.Time, fields []Field)
	Clone() Encoder
	AddFields([]Field)
}
