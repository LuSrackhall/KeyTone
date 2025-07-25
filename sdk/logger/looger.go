/*
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 最大日志大小(MB)
const LOG_MAX_SIZE = 4

// 日志保留天数
const LOG_MAX_DAYS = 7

// 通过环境变量SDK_MODE来判断, 输出到 标准输出 还是 文件
var Logger *slog.Logger
var LoggerDebug *slog.Logger
var ProgramLevel *slog.LevelVar

func InitLogger(log_file_pathAndName string) {
	ProgramLevel = new(slog.LevelVar)
	var handler slog.Handler
	var handlerDebug slog.Handler

	// 日志轮转功能：检查文件大小并执行轮转
	if shouldRotateLog(log_file_pathAndName) {
		if err := rotateLogFile(log_file_pathAndName); err != nil {
			slog.Error("Failed to rotate log file", "err", err)
		}
	}

	// 清理过期日志
	if err := cleanOldLogs(log_file_pathAndName); err != nil {
		slog.Error("Failed to clean old logs", "err", err)
	}

	// 打开一个文件用于写入日志
	// file, err := os.OpenFile("log.jsonl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	file, err := os.OpenFile(log_file_pathAndName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "err", err)
		os.Exit(1)
	}

	// 根据环境变量的值来决定日志的输出位置
	// (debug release两个选项可选,默认只要不显式指定相关环境变量为debug, 就默认为release模式--即release是默认不需要显式指定的)
	if os.Getenv("SDK_MODE") == "debug" {
		// 开发环境(即debug)：日志输出到标准输出和文件
		// handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: ProgramLevel})
		handlerDebug = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: ProgramLevel})
		LoggerDebug = slog.New(handlerDebug)
	} else {
		// 生产环境(即release)：日志只输出到文件
		// handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: ProgramLevel})
		// handlerDebug = nil
		LoggerDebug = nil
	}

	handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: ProgramLevel})
	Logger = slog.New(handler)
}

/**
 * 对于调试过程起到推荐作用的详细信息, 可使用以下自定义的Log函数
 */

// 自定义的 Debug 函数
func Debug(msg string, keyvals ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keyvals = append([]interface{}{"file", file, "line", line}, keyvals...)
	Logger.Debug(msg, keyvals...)
	if LoggerDebug != nil {
		LoggerDebug.Debug(msg, keyvals...)
	}
}

// 自定义的 Info 函数
func Info(msg string, keyvals ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keyvals = append([]interface{}{"file", file, "line", line}, keyvals...)
	Logger.Info(msg, keyvals...)
	if LoggerDebug != nil {
		LoggerDebug.Info(msg, keyvals...)
	}
}

// 自定义的 Warn 函数
func Warn(msg string, keyvals ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keyvals = append([]interface{}{"file", file, "line", line}, keyvals...)
	Logger.Warn(msg, keyvals...)
	if LoggerDebug != nil {
		LoggerDebug.Warn(msg, keyvals...)
	}
}

// 自定义的 Error 函数
func Error(msg string, keyvals ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keyvals = append([]interface{}{"file", file, "line", line}, keyvals...)
	Logger.Error(msg, keyvals...)
	if LoggerDebug != nil {
		LoggerDebug.Error(msg, keyvals...)
	}
}

// 检查是否需要轮转日志文件
func shouldRotateLog(filePath string) bool {

	maxSize := LOG_MAX_SIZE * 1024 * 1024

	if maxSize == 0 {
		// 默认最大日志大小10MB
		maxSize = 10 * 1024 * 1024 // 10MB in bytes
	}

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		slog.Error("Failed to get log file info", "err", err)
		return false
	}

	return fileInfo.Size() > int64(maxSize)
}

// 执行日志文件轮转
func rotateLogFile(filePath string) error {
	// 生成带时间戳的新文件名
	dir, filename := filepath.Split(filePath)
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("20060102-150405")
	newFilename := fmt.Sprintf("%s-%s%s", base, timestamp, ext)
	newPath := filepath.Join(dir, newFilename)

	// 重命名当前日志文件
	if err := os.Rename(filePath, newPath); err != nil {
		return err
	}
	return nil
}

// 清理过期日志
func cleanOldLogs(currentLogPath string) error {
	retentionDays := LOG_MAX_DAYS
	if retentionDays == 0 { // Fallback if LOG_MAX_DAYS is explicitly 0, use a reasonable default
		// 默认保留7天 (Using the constant LOG_MAX_DAYS directly)
		retentionDays = 7
	}

	// 获取日志目录和文件名
	logDir, filename := filepath.Split(currentLogPath)
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	for _, file := range files {
		filePath := filepath.Join(logDir, file.Name())

		// Skip the current log file
		if filePath == currentLogPath {
			continue
		}

		// Check filename pattern (base-timestamp.ext)
		// Extract timestamp from filename
		fileName := file.Name()
		if !strings.HasPrefix(fileName, base+"-") {
			continue
		}

		// Expected format: base-YYYYMMDD-HHmmss.ext
		// Example: mylog-20240725-103000.jsonl
		parts := strings.Split(strings.TrimSuffix(fileName, filepath.Ext(fileName)), "-")
		if len(parts) < 3 { // Ensure there are enough parts for date and time
			continue
		}

		// Attempt to parse the timestamp from the filename
		timestampStr := parts[len(parts)-2] + "-" + parts[len(parts)-1] // e.g., "20060102-150405"

		// Adjust the format string to match the format you used for generating the timestamp
		logTime, err := time.Parse("20060102-150405", timestampStr)
		if err != nil {
			slog.Warn("Could not parse timestamp from log file name, skipping cleanup for this file", "file", fileName, "err", err)
			continue
		}

		// Delete expired files based on the timestamp in the filename
		if logTime.Before(cutoffTime) {
			if err := os.Remove(filePath); err != nil {
				slog.Error("Failed to remove old log file", "file", fileName, "err", err)
			} else {
				slog.Info("Removed old log file", "file", fileName)
			}
		}
	}
	return nil
}
