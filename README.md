# Riky

**Important: 项目仍在开发中，功能不完善**

一个高性能 Go 结构化日志库，参考 Zap 设计，提供零分配 (zero-allocation) JSON 编码、強类型字段和动态级别支持。

## 特性

- **零分配编码**: 手写 JSON 序列化，避免 `encoding/json` 反射开销
- **强类型字段**: 支持 `Int`, `String`, `Float`, `Bool`, `Time`, `Duration`, `Error`
- **Buffer 池化**: 使用 `sync.Pool` 复用缓冲区，减少 GC
- **动态级别**: 运行时通过 `AtomicLevel` 切换日志级别，无需重启
- **全局 Logger**: 简洁的全局 API (`logger.Info()`)
- **上下文携带**: `With()` 支持链式调用

## 快速开始

```go
package main

import (
    "github.com/Yonagi04/Riky/logger"
)

func main() {
    // 基础用法
    log := logger.NewLogger()
    log.Info("user login", logger.Int("uid", 12345))

    // 上下文携带
    reqLog := log.With(logger.String("request_id", "abc-123"))
    reqLog.Info("processing request")

    // 动态级别
    atomicLevel := logger.NewAtomicLevelAt(logger.DebugLevel)
    dynamicLog := logger.NewLoggerWithAtomicLevel(atomicLevel)

    dynamicLog.Debug("debug message")
    atomicLevel.SetLevel(logger.WarnLevel) // 运行时切换
    dynamicLog.Info("this will be filtered")

    // 全局 Logger
    logger.Info("global logger message")
    logger.SetLevel(logger.DebugLevel)
    logger.Debug("global debug")
}
```

输出示例:

```json
{"level":"INFO","msg":"user login","uid":12345}
{"level":"INFO","msg":"processing request","request_id":"abc-123"}
```

## 安装

```bash
go get github.com/Yonagi04/Riky
```

## API 参考

### 创建 Logger

```go
// 默认 (输出到 stdout)
log := logger.NewLogger()

// 自定义输出
log := logger.NewLoggerWithWriter(w)

// 动态级别
log := logger.NewLoggerWithAtomicLevel(level)
```

### 日志方法

```go
log.Debug(msg string, fields...)
log.Info(msg string, fields...)
log.Warn(msg string, fields...)
log.Error(msg string, fields...)
log.Panic(msg string, fields...)
log.Fatal(msg string, fields...)
```

### 字段构造

```go
logger.Int(key string, value int)
logger.String(key string, value string)
logger.Float(key string, value float64)
logger.Bool(key string, value bool)
logger.Time(key string, value time.Time)
logger.Duration(key string, value time.Duration)
logger.Error(key string, value error)
```

### 动态级别

```go
level := logger.NewAtomicLevelAt(logger.DebugLevel)
level.SetLevel(logger.InfoLevel)

// 或通过 Logger
log.SetLevel(logger.WarnLevel)
currentLevel := log.Level()
```

## 项目结构

```
Riky/
├── logger/           # 主包
│   ├── logger.go    # 用户 API
│   ├── core.go    # 核心逻辑
│   ├── level.go   # 日志级别
│   ├── field.go   # 强类型字段
│   ├── buffer.go # Buffer 池化
│   └── atomic_level.go # 动态级别
├── encoder/         # 编码器
│   ├── encoder.go # 接口定义
│   └── json.go    # JSON 实现
└── example/       # 使用示例
```

## 性能

参考 Zap 的设计理念:

- 手写 `append` 方式序列化 JSON，避免反射
- Buffer 池化复用内存
- 懒编码: 只有级别通过时才序列化
- 零分配 (zero-allocation) 目标

## 演进路线

| 阶段 | 目标 |
|------|------|
| MVP | 基础日志、强类型字段、JSON 编码 |
| Performance | 零分配、Buffer 池化、懒编码 |
| Production | 异步写入、Hooks、Caller、Stacktrace |

## 许可证

MIT