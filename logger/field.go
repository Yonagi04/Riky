package logger

import (
	"math"
	"time"

	"github.com/Yonagi04/Riky/encoder"
)

type Field = encoder.Field

func Int(key string, value int) Field {
	return Field{
		Key:     key,
		Type:    encoder.Int64Type,
		Integer: int64(value),
	}
}

func String(key string, value string) Field {
	return Field{
		Key:    key,
		Type:   encoder.StringType,
		String: value,
	}
}

func Float(key string, value float64) Field {
	return Field{
		Key:   key,
		Type:  encoder.Float64Type,
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
		Type:    encoder.BoolType,
		Integer: intValue,
	}
}

func Time(key string, value time.Time) Field {
	return Field{
		Key:     key,
		Type:    encoder.TimeType,
		Integer: value.UnixNano(),
	}
}

func Duration(key string, value time.Duration) Field {
	return Field{
		Key:     key,
		Type:    encoder.DurationType,
		Integer: value.Nanoseconds(),
	}
}

func Error(key string, value error) Field {
	if value != nil {
		return Field{
			Key:    key,
			Type:   encoder.ErrorType,
			String: value.Error(),
		}
	}
	return Field{}
}
