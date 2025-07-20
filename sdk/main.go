/**
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

package main

import (
	audioPackageConfig "KeyTone/audioPackage/config"
	"KeyTone/config"
	"KeyTone/keyEvent"
	"KeyTone/logger"
	"KeyTone/server"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 定义配置文件路径的命令行参数
var ConfigPath string

// 定义音频包根目录的路径
// var AudioPackagePath string  // 定义位置不应该在此处, 已将其移动到音频包模块中去

// 定义日志文件路径的命令行参数
var LogPathAndName string

// https://github.com/gopxl/beep/issues/179 此测试代码块对于简单的内存泄漏检测很有帮助, 之前曾借助其定位过beep的内存泄漏问题。
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))             // 当前堆上分配的内存（MiB）。表示程序当前正在使用的堆内存量。
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)) // 累计分配的内存总量（MiB）。这个值只增不减，但如果程序持续运行且没有内存泄漏，它的增长速度应该会变慢（稳定）。
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))               // 从系统获得的总内存（MiB）。这个值包括堆、栈和其他系统分配的内存。
	fmt.Printf("\tNumGC = %v\n", m.NumGC)                    // 垃圾回收的次数。通过这个可以看GC是否在运行。
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func init() {
	go func() {
		for {
			PrintMemUsage()
			time.Sleep(5 * time.Second)
		}
	}()
	// 设置环境变量
	err := godotenv.Load()
	// 在 InitLogger 之前, 使用slog的log信息(与我们的logger模块无关), 此类信息仅在终端上看即可, 输出到日志中也无意义。
	if err != nil {
		// 没必要因为这个err退出程序, .env文件在本项目中, 主要作为开发文件使用。(后面真要上配置文件的话, 也是使用.json格式的)
		slog.Warn("无法加载.env文件", "err", err)
	} else {
		slog.Info(".env文件已被正确加载", "SDK_MODE", os.Getenv("SDK_MODE"))
	}

	// 获取命令行参数中的传入值
	{

		// 如果路径不存在, 则使用当前目录作为路径
		// * 第一个参数是指向一个字符串变量的指针，用于存储解析后的值。
		// * 第二个参数是命令行参数的名称（在命令行中使用）。  用户在使用时 go run main.go -configPath=./path
		// * 第三个参数是默认值（如果用户没有提供这个参数，则使用默认值）。
		// * 第四个参数是这个参数的描述（帮助信息）。
		flag.StringVar(&ConfigPath, "configPath", ".", "Path to the config file")
		flag.StringVar(&audioPackageConfig.AudioPackagePath, "audioPackagePath", "./temporaryDebug", "Path to the Audio Package Root Dir")
		flag.StringVar(&LogPathAndName, "logPathAndName", "./log.jsonl", "Path and name to the log file")

		// 解析命令行参数
		flag.Parse()

		// 使用命令行参数
		// ...
		slog.Info("命行参数已正确解析", "configPath", ConfigPath, "audioPackagePath", audioPackageConfig.AudioPackagePath, "logPathAndName", LogPathAndName)
	}

	// 初始化模块
	{

		// 初始化日志模块(并顺便初始化gin的MODE), 主要是为了输出到日志中, 便于在用户使用过程中记录bug数据。
		{
			logger.InitLogger(LogPathAndName)

			// 设置日志级别(此处主要用于开发过程中, 自己可随时进行调整的级别设置)
			logger.ProgramLevel.Set(slog.LevelDebug)
			// logger.ProgramLevel.Set(slog.LevelInfo)
			// logger.ProgramLevel.Set(slog.LevelWarn)
			// logger.ProgramLevel.Set(slog.LevelError)

			if os.Getenv("SDK_MODE") != "debug" {
				// 设置log库, 在正式release中的默认级别
				logger.ProgramLevel.Set(slog.LevelInfo)

				// 设置 gin 框架, 在正式release中的 MODE 为 "release"
				gin.SetMode(gin.ReleaseMode)
			}

			logger.Info("日志模块已开始正常运行, Getenv值已获取。 ", "SDK_MODE", os.Getenv("SDK_MODE"), "GIN_MODE", os.Getenv("GIN_MODE"))
		}

		// 初始化配置模块
		{
			// 检查指定的路径是否存在
			if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
				// 如果路径不存在，创建路径
				err := os.MkdirAll(ConfigPath, os.ModePerm)
				if err != nil {
					logger.Error("配置文件路径创建时出错。", "err", err.Error())
				} else {
					logger.Info("配置文件路径创建成功。", "你的配置文件路径为", ConfigPath)
				}
			} else if err != nil {
				logger.Error("检查配置文件路径时出错。", "err", err.Error())
			} else {
				logger.Info("配置文件路径已存在且无异常。", "你的配置文件路径为", ConfigPath)
			}
			config.ConfigRun(ConfigPath)
		}

	}

}

func main() {
	go server.ServerRun()
	keyEvent.KeyEventListen()
}
