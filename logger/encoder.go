package logger

import (
	"fmt"
	"math"
	"strconv"
)

// hexDigits 用于将字节转换为十六进制字符
var hexDigits = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

type Encoder interface {
	Encode(buf *Buffer, msg string, level Level, fields []Field)
	Clone() Encoder
	AddFields([]Field)
}

type jsonEncoder struct {
	buf []byte
}

func (e *jsonEncoder) Encode(buf *Buffer, msg string, level Level, fields []Field) {
	buf.AppendString(`{"level":"`)
	buf.AppendString(level.String())
	buf.AppendString(`","msg":"`)
	e.safeAppendString(buf, msg)
	buf.AppendByte('"')
	buf.AppendBytes(e.buf) // 输出 From With() 累积的持久化字段

	// 使用临时 buffer 处理本次调用传入的临时字段，不修改 e.buf
	tempBuf := getBuffer()
	defer putBuffer(tempBuf)
	e.addFieldsToBuffer(tempBuf, fields)
	buf.AppendBytes(tempBuf.bs)

	buf.AppendString("}\n")
}

// addFieldsToBuffer 将字段添加到指定的 buffer，不修改 encoder 内部状态
func (e *jsonEncoder) addFieldsToBuffer(buf *Buffer, fields []Field) {
	for _, f := range fields {
		e.addKeyToBuffer(buf, f.Key)
		switch f.Type {
		case StringType:
			buf.AppendString("\"")
			e.safeAppendString(buf, f.String)
			buf.AppendString("\"")
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
			if f.Interface != nil {
				buf.AppendString(f.Interface.(error).Error())
			}
		case ReflectType:
			if f.Interface != nil {
				buf.AppendString(fmt.Sprintf("%+v", f.Interface))
			}
		}
	}
}

func (e *jsonEncoder) addKeyToBuffer(buf *Buffer, key string) {
	buf.AppendByte(',')
	buf.AppendString("\"")
	buf.AppendString(key)
	buf.AppendString("\":")
}

func (e *jsonEncoder) appendString(s string) {
	e.buf = append(e.buf, '"')
	e.buf = append(e.buf, s...)
	e.buf = append(e.buf, '"')
}

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
			if f.Interface != nil {
				e.appendString(f.Interface.(error).Error())
			}
		case ReflectType:
			if f.Interface != nil {
				e.appendString(fmt.Sprintf("%+v", f.Interface))
			}
		default:
			panic("unhandled default case")
		}
	}
}

func (e *jsonEncoder) addKey(key string) {
	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
}

func (e *jsonEncoder) Clone() Encoder {
	newBuf := make([]byte, len(e.buf))
	copy(newBuf, e.buf)
	return &jsonEncoder{buf: newBuf}
}
