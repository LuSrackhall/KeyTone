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
	"log/slog"
	"os"
	"runtime"
)

// 通过环境变量SDK_MODE来判断, 输出到 标准输出 还是 文件
var Logger *slog.Logger
var LoggerDebug *slog.Logger
var ProgramLevel *slog.LevelVar

func InitLogger(log_file_pathAndName string) {
	ProgramLevel = new(slog.LevelVar)
	var handler slog.Handler
	var handlerDebug slog.Handler

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
