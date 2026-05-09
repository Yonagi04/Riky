package logger

import (
	"fmt"
	"strings"
)

// Level 表示日志级别
type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// LevelEnabler 判断给定级别是否应该输出日志
type LevelEnabler interface {
	Enabled(level Level) bool
}

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return "INFO"
	}
}

// Enabled 实现 LevelEnabler 接口
// 当 target >= l 时返回 true
func (l Level) Enabled(target Level) bool {
	return target >= l
}

// ParseLevel 将字符串解析为 Level
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown level: %s", s)
	}
}
