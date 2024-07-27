package config

import (
	"KeyTone/logger"

	"github.com/spf13/viper"
)

// Setting.json is the configuration for the application
func ConfigRun(path string) {
	// 设置配置文件名称和类型
	viper.SetConfigName("KeyToneSetting")
	viper.SetConfigType("json")

	// 添加配置文件路径
	viper.AddConfigPath(path)
	// viper.AddConfigPath(path2) // 越靠下, 优先级越高

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，创建默认设置的配置文件
			logger.Info("未找到默认的配置文件，正在创建...")
			createDefaultSetting()
		} else {
			// 其他错误
			logger.Error("读取设置文件时发生致命错误", "err", err.Error())
		}
	}

	// 使用配置
	value := viper.GetString("key")
	logger.Debug("打印测试用的key值", "value", value)

	// 设置新配置值, 并将设置的值保存到配置文件
	viper.Set("key", "newValue")
	if err := viper.WriteConfig(); err != nil {
		logger.Error("向配置文件保存设置时发生致命错误", "err", err.Error())
	}
}

// 创建默认设置的配置文件
func createDefaultSetting() {
	viper.Set("key", "defaultValue")
	if err := viper.SafeWriteConfig(); err != nil {
		logger.Error("创建默认设置文件时发生致命错误", "err", err.Error())

	}
}
