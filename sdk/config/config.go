package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Setting.json is the configuration for the application
func ConfigRun(path string) {
	// 设置配置文件名称和类型
	viper.SetConfigName("setting")
	viper.SetConfigType("json")

	// 添加配置文件路径
	viper.AddConfigPath(path)
	// viper.AddConfigPath(path2) // 越靠下, 优先级越高

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，创建默认设置的配置文件
			fmt.Println("未找到设置文件，正在创建...")
			createDefaultSetting()
		} else {
			// 其他错误
			fmt.Printf("读取设置文件时出错: %s", err)
		}
	}

	// 使用配置
	value := viper.GetString("key")
	fmt.Println("Value:", value)

	// 设置新配置值, 并将设置的值保存到配置文件
	viper.Set("key", "newValue")
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("写入设置时出错, 错误为: %s", err)
	}
}

// 创建默认设置的配置文件
func createDefaultSetting() {
	viper.Set("key", "defaultValue")
	if err := viper.SafeWriteConfig(); err != nil {
		fmt.Printf("创建默认设置文件时出错, 错误为: %s", err)
	}
}
