package logger

import (
	"math"
	"time"
)

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

func Int(key string, value int) Field {
	return Field{
		Key:     key,
		Type:    Int64Type,
		Integer: int64(value),
	}
}

func String(key string, value string) Field {
	return Field{
		Key:    key,
		Type:   StringType,
		String: value,
	}
}

func Float(key string, value float64) Field {
	return Field{
		Key:   key,
		Type:  Float64Type,
		Float: math.Float64bits(value),
	}
}

func Bool(key string, value bool) Field {
	var intValue int64
	if value {
		intValue = 1
	}
	return Field{
		Key:     key,
		Type:    BoolType,
		Integer: intValue,
	}
}

func Time(key string, value time.Time) Field {
	return Field{
		Key:     key,
		Type:    TimeType,
		Integer: value.UnixNano(),
	}
}

func Duration(key string, value time.Duration) Field {
	return Field{
		Key:     key,
		Type:    DurationType,
		Integer: value.Nanoseconds(),
	}
}

func Error(key string, value error) Field {
	if value != nil {
		return Field{
			Key:    key,
			Type:   ErrorType,
			String: value.Error(),
		}
	}
	return Field{}
}
