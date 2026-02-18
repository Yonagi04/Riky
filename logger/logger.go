package logger

import (
	"os"
	"sync/atomic"
)

type Logger struct {
	core *Core
}

func NewLogger() *Logger {
	core := NewCore(os.Stdout, NewEncoder(), InfoLevel)
	return &Logger{core: core}
}

// NewLoggerWithAtomicLevel 创建带动态级别的 Logger
func NewLoggerWithAtomicLevel(level *AtomicLevel) *Logger {
	core := NewCoreWithEnabler(os.Stdout, NewEncoder(), level)
	return &Logger{core: core}
}

// SetLevel 设置 Logger 的日志级别
// 注意：仅当 Logger 使用 AtomicLevel 时才能动态生效
func (l *Logger) SetLevel(level Level) {
	if al, ok := l.core.levelEnabler.(*AtomicLevel); ok {
		al.SetLevel(level)
	}
}

// Level 获取 Logger 当前的日志级别
func (l *Logger) Level() Level {
	switch enabler := l.core.levelEnabler.(type) {
	case *AtomicLevel:
		return enabler.Level()
	case Level:
		return enabler
	default:
		return InfoLevel
	}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	err := l.core.Write(DebugLevel, msg, fields)
	if err != nil {
		return
	}
}

func (l *Logger) Info(msg string, fields ...Field) {
	err := l.core.Write(InfoLevel, msg, fields)
	if err != nil {
		return
	}
}

func (l *Logger) Warn(msg string, fields ...Field) {
	err := l.core.Write(WarnLevel, msg, fields)
	if err != nil {
		return
	}
}

func (l *Logger) Error(msg string, fields ...Field) {
	err := l.core.Write(ErrorLevel, msg, fields)
	if err != nil {
		return
	}
}

func (l *Logger) Panic(msg string, fields ...Field) {
	err := l.core.Write(PanicLevel, msg, fields)
	if err != nil {
		return
	}
	panic(msg)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	err := l.core.Write(FatalLevel, msg, fields)
	if err != nil {
		return
	}
	os.Exit(1)
}

func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		core: l.core.With(fields),
	}
}

// ========================================
// 全局 Logger 支持
// ========================================

var defaultLogger atomic.Pointer[Logger]

func init() {
	// 初始化默认 Logger，使用 AtomicLevel 以支持动态级别
	level := NewAtomicLevelAt(InfoLevel)
	l := NewLoggerWithAtomicLevel(&level)
	defaultLogger.Store(l)
}

// SetDefault 设置全局 Logger
func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

// Default 获取全局 Logger
func Default() *Logger {
	return defaultLogger.Load()
}

// SetLevel 设置全局 Logger 的日志级别
func SetLevel(level Level) {
	Default().SetLevel(level)
}

// GetLevel 获取全局 Logger 的日志级别
func GetLevel() Level {
	return Default().Level()
}

// Debug 使用全局 Logger 输出 Debug 级别日志
func Debug(msg string, fields ...Field) {
	Default().Debug(msg, fields...)
}

// Info 使用全局 Logger 输出 Info 级别日志
func Info(msg string, fields ...Field) {
	Default().Info(msg, fields...)
}

// Warn 使用全局 Logger 输出 Warn 级别日志
func Warn(msg string, fields ...Field) {
	Default().Warn(msg, fields...)
}

// LogError 使用全局 Logger 输出 Error 级别日志
// 注意：命名为 LogError 以避免与 field.go 中的 Error 函数冲突
func LogError(msg string, fields ...Field) {
	Default().Error(msg, fields...)
}

// Panic 使用全局 Logger 输出 Panic 级别日志
func Panic(msg string, fields ...Field) {
	Default().Panic(msg, fields...)
}

// Fatal 使用全局 Logger 输出 Fatal 级别日志
func Fatal(msg string, fields ...Field) {
	Default().Fatal(msg, fields...)
}
