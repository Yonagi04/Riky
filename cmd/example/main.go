package main

import (
	"github.com/Yonagi04/Riky/riky"
)

func main() {
	// ========================================
	// 示例 1: 基础用法（静态级别）
	// ========================================
	log := riky.NewLogger()
	log.Info("user login", riky.Int("uid", 12345), riky.String("status", "success"))

	// 派生用法 (Contextual logging)
	subLogger := log.With(riky.String("request_id", "abc-123"))
	subLogger.Info("processing request")

	log.Info("another log without request_id")

	// ========================================
	// 示例 2: 动态级别 (AtomicLevel)
	// ========================================
	atomicLevel := riky.NewAtomicLevelAt(riky.DebugLevel)
	dynamicLog := riky.NewLoggerWithAtomicLevel(atomicLevel)

	// 初始为 Debug 级别，所有日志都会输出
	dynamicLog.Debug("debug message - will be printed")
	dynamicLog.Info("info message - will be printed")

	// 动态切换到 Info 级别
	atomicLevel.SetLevel(riky.InfoLevel)
	dynamicLog.Debug("debug message - will NOT be printed") // 被过滤
	dynamicLog.Info("info message after level change")      // 正常输出

	// 通过 Logger 方法设置级别
	dynamicLog.SetLevel(riky.WarnLevel)
	dynamicLog.Info("info message - will NOT be printed") // 被过滤
	dynamicLog.Warn("warn message - will be printed")     // 正常输出

	// ========================================
	// 示例 3: 全局 Logger
	// ========================================
	riky.Info("global logger message")
	riky.SetLevel(riky.DebugLevel)
	riky.Debug("global debug after level change")
	riky.Info("current level", riky.String("level", riky.GetLevel().String()))
}
