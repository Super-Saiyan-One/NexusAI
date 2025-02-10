package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"nexus-ai/constant"

	"github.com/gin-gonic/gin"
)

const (
	logINFO  = "INFO"
	logWARN  = "WARN"
	logERROR = "ERROR"

	colorReset = "\033[0m"

	// 普通日志的颜色代码
	colorRed    = "\033[91m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"

	// 系统日志的颜色代码（使用更鲜艳的颜色）
	sysColorRed    = "\033[31m" // 亮红色
	sysColorGreen  = "\033[32m" // 亮绿色
	sysColorYellow = "\033[33m" // 亮黄色
)

var maxLogCount = constant.LogMaxCount
var logDir = constant.LogDir

var logCount int
var setupLogWorking bool
var setupLogLock sync.Mutex

func SetupLog() {
	if logDir != "" {
		ok := setupLogLock.TryLock()
		if !ok {
			log.Println("SetupLog is working")
			return
		}
		logPath := filepath.Join(logDir, fmt.Sprintf("nexus-ai-%s.log", time.Now().Format("2006-01-02 15-04-05")))
		fd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("SetupLog failed to open log file:", err)
		}
		gin.DefaultWriter = io.MultiWriter(os.Stdout, fd)
		gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, fd)
		defer func() {
			setupLogLock.Unlock()
			setupLogWorking = false
			SysInfo("SetupLog Successfully")
		}()
	}
}

func logCommon(ctx context.Context, logType string, logContent string) {
	writer := gin.DefaultErrorWriter // 默认输出到错误输出
	if logType == logINFO {          // 日志类型为INFO时，输出到标准输出
		writer = gin.DefaultWriter
	}
	userID := GetContextValue(ctx.Value(constant.UserIDKey), "user less")
	requestID := GetContextValue(ctx.Value(constant.RequestIDKey), "test_request_id")
	t := time.Now().Format("2006-01-02 15:04:05")
	switch logType {
	case logINFO:
		_, _ = fmt.Fprintf(writer, "%s[%s]%s    %v | %s | %s | %s%s%s \n", colorGreen, logType, colorReset, t, userID, requestID, colorGreen, logContent, colorReset)
	case logWARN:
		_, _ = fmt.Fprintf(writer, "%s[%s]%s    %v | %s | %s | %s%s%s \n", colorYellow, logType, colorReset, t, userID, requestID, colorYellow, logContent, colorReset)
	case logERROR:
		_, _ = fmt.Fprintf(writer, "%s[%s]%s   %v | %s | %s | %s%s%s \n", colorRed, logType, colorReset, t, userID, requestID, colorRed, logContent, colorReset)
	}
	logCount++
	if logCount > maxLogCount && !setupLogWorking { // 如果日志数量超过最大值，并且没有正在设置日志，则设置日志
		logCount = 0
		setupLogWorking = true
		go func() {
			SetupLog()
		}()
	}
}

func LogInfo(ctx context.Context, logContent string) { // 信息日志
	logCommon(ctx, logINFO, logContent)
}

func LogWarn(ctx context.Context, logContent string) { // 警告日志
	logCommon(ctx, logWARN, logContent)
}

func LogError(ctx context.Context, logContent string) { // 错误日志
	logCommon(ctx, logERROR, logContent)
}

func SysInfo(s string) { // 系统信息日志
	t := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(gin.DefaultWriter, "%s[SYSINFO]%s %v | %s%s%s \n", sysColorGreen, colorReset, t, sysColorGreen, s, colorReset)
}

func SysWarn(s string) { // 系统警告日志
	t := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(gin.DefaultWriter, "%s[SYSWARN]%s %v | %s%s%s \n", sysColorYellow, colorReset, t, sysColorYellow, s, colorReset)
}

func SysError(s string) { // 系统错误日志
	t := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "%s[SYSERR]%s  %v | %s%s%s \n", sysColorRed, colorReset, t, sysColorRed, s, colorReset)
}

func FatalLog(s string) { // 致命错误日志 系统退出
	t := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(gin.DefaultWriter, "%s[SYSFAIL]%s %v | %s%s%s \n", sysColorRed, colorReset, t, sysColorRed, s, colorReset)
	os.Exit(1)
}
