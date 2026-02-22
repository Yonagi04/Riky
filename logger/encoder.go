package logger

import (
	"math"
	"strconv"
	"time"
)

// hexDigits 用于将字节转换为十六进制字符
var hexDigits = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

// JSON 字段名常量
// 集中管理，便于维护和修改，避免硬编码字符串分散在代码中
const (
	jsonFieldLevel = "level"
	jsonFieldMsg   = "msg"
	jsonFieldTime  = "time"
)

// Encoder 定义日志编码器接口
type Encoder interface {
	Encode(buf *Buffer, msg string, level Level, time time.Time, fields []Field)
	Clone() Encoder
	AddFields([]Field)
}

// jsonEncoder JSON 格式的日志编码器
// buf 存储 With() 累积的持久化字段，跨多次 Encode 调用保持数据
type jsonEncoder struct {
	buf []byte
}

// NewEncoder 创建新的 JSON 编码器
// 预分配 buf 容量，避免后续 AddFields 时频繁扩容
func NewEncoder() *jsonEncoder {
	return &jsonEncoder{
		buf: make([]byte, 0, 256),
	}
}

// Encode 将日志编码为 JSON 格式
// 利用 Go 编译器优化：append(dst, src...) 会被优化为高效的 memmove
func (e *jsonEncoder) Encode(buf *Buffer, msg string, level Level, t time.Time, fields []Field) {
	// 开始 JSON 对象
	buf.AppendByte('{')

	// level 字段（第一个字段，无前导逗号）
	e.writeFirstStringField(buf, jsonFieldLevel, level.String())

	// msg 字段
	e.writeStringField(buf, jsonFieldMsg, msg)

	// time 字段
	e.writeTimeField(buf, jsonFieldTime, t)

	// From With() 累积的持久化字段
	// Go 编译器会将 append 优化为高效的 memmove
	buf.bs = append(buf.bs, e.buf...)

	// 本次调用传入的临时字段（直接写入主 buffer，避免额外分配）
	e.addFieldsToBuffer(buf, fields)

	// 结束 JSON 对象
	buf.AppendString("}\n")
}

// writeFirstStringField 写入第一个字符串字段（无前导逗号）
func (e *jsonEncoder) writeFirstStringField(buf *Buffer, key, value string) {
	buf.AppendByte('"')
	buf.AppendString(key)
	buf.AppendString("\":\"")
	e.safeAppendString(buf, value)
	buf.AppendByte('"')
}

// writeStringField 写入字符串字段（带前导逗号）
func (e *jsonEncoder) writeStringField(buf *Buffer, key, value string) {
	buf.AppendString(",\"")
	buf.AppendString(key)
	buf.AppendString("\":\"")
	e.safeAppendString(buf, value)
	buf.AppendByte('"')
}

// writeTimeField 写入时间字段
func (e *jsonEncoder) writeTimeField(buf *Buffer, key string, t time.Time) {
	buf.AppendString(",\"")
	buf.AppendString(key)
	buf.AppendString("\":\"")

	year, month, day := t.Date()
	hour, minute, sec := t.Clock()

	// 日期部分: YYYY-MM-DD
	buf.AppendInt(int64(year))
	buf.AppendByte('-')
	if month < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int64(month))
	buf.AppendByte('-')
	if day < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int64(day))

	// 时间部分: HH:MM:SS
	buf.AppendByte(' ')
	if hour < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int64(hour))
	buf.AppendByte(':')
	if minute < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int64(minute))
	buf.AppendByte(':')
	if sec < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int64(sec))

	buf.AppendByte('"')
}

// addFieldsToBuffer 将字段添加到指定的 buffer，不修改 encoder 内部状态
func (e *jsonEncoder) addFieldsToBuffer(buf *Buffer, fields []Field) {
	for _, f := range fields {
		e.addKeyToBuffer(buf, f.Key)
		switch f.Type {
		case StringType:
			buf.AppendByte('"')
			e.safeAppendString(buf, f.String)
			buf.AppendByte('"')
		case Int64Type:
			buf.AppendInt(f.Integer)
		case Float64Type:
			buf.bs = strconv.AppendFloat(buf.bs, math.Float64frombits(f.Float), 'f', -1, 64)
		case BoolType:
			if f.Integer == 1 {
				buf.AppendString("true")
			} else {
				buf.AppendString("false")
			}
		case TimeType:
			buf.AppendInt(f.Integer)
		case DurationType:
			buf.AppendInt(f.Integer)
		case ErrorType:
			buf.AppendByte('"')
			e.safeAppendString(buf, f.String)
			buf.AppendByte('"')
		default:
			panic("unhandled default case")
		}
	}
}

// addKeyToBuffer 向 buffer 添加 JSON 键名
func (e *jsonEncoder) addKeyToBuffer(buf *Buffer, key string) {
	buf.AppendString(",\"")
	buf.AppendString(key)
	buf.AppendString("\":")
}

// safeAppendString 安全地追加字符串，处理 JSON 特殊字符转义
func (e *jsonEncoder) safeAppendString(buf *Buffer, s string) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			buf.AppendString(`\"`)
		case '\\':
			buf.AppendString(`\\`)
		case '\n':
			buf.AppendString(`\n`)
		case '\r':
			buf.AppendString(`\r`)
		case '\t':
			buf.AppendString(`\t`)
		case '\b':
			buf.AppendString(`\b`)
		case '\f':
			buf.AppendString(`\f`)
		default:
			// 处理不可见控制字符 (JSON 要求转义为 \u00xx)
			if c < 0x20 {
				buf.AppendString(`\u00`)
				buf.AppendByte(hexDigits[(c>>4)&0x0F])
				buf.AppendByte(hexDigits[c&0x0F])
			} else {
				buf.AppendByte(c)
			}
		}
	}
}

// AddFields 添加持久化字段到 encoder
func (e *jsonEncoder) AddFields(fields []Field) {
	for _, f := range fields {
		e.addKey(f.Key)
		switch f.Type {
		case StringType:
			e.appendString(f.String)
		case Int64Type:
			e.buf = strconv.AppendInt(e.buf, f.Integer, 10)
		case Float64Type:
			e.buf = strconv.AppendFloat(e.buf, math.Float64frombits(f.Float), 'f', -1, 64)
		case BoolType:
			if f.Integer == 1 {
				e.buf = append(e.buf, "true"...)
			} else {
				e.buf = append(e.buf, "false"...)
			}
		case TimeType:
			e.buf = strconv.AppendInt(e.buf, f.Integer, 10)
		case DurationType:
			e.buf = strconv.AppendInt(e.buf, f.Integer, 10)
		case ErrorType:
			e.appendString(f.String)
		default:
			panic("unhandled default case")
		}
	}
}

// addKey 向 encoder 内部 buf 添加 JSON 键名
func (e *jsonEncoder) addKey(key string) {
	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
}

// appendString 向 encoder 内部 buf 添加字符串值
func (e *jsonEncoder) appendString(s string) {
	e.buf = append(e.buf, '"')
	e.buf = append(e.buf, s...)
	e.buf = append(e.buf, '"')
}

// Clone 克隆 encoder，深拷贝 buf
func (e *jsonEncoder) Clone() Encoder {
	newBuf := make([]byte, len(e.buf))
	copy(newBuf, e.buf)
	return &jsonEncoder{buf: newBuf}
}
