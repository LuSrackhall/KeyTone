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
			createDefaultConfig()
			logger.Info("配置文件创建并载入成功")
		} else {
			// 其他错误
			logger.Error("读取设置文件时发生致命错误", "err", err.Error())
		}
	} else {
		logger.Info("配置文件已加载, 正在与DefaultConfig进行diff和增量载入...")
		diffAndUpdateDefaultConfig()
		logger.Info("配置文件diff和增量载入完成")
	}
}

// 恢复默认设置(删除掉已有配置文件后, 再调用“createDefaultConfig()"函数就好了)
func RestoreDefaultSetting() {

}

// 获取配置
// * 如果想获取这个json文件的所有配置, 则可用key="get_all_value"作为键值来查询
func GetValue(key string) any {
	if key == "get_all_value" {
		return viper.AllSettings()
	} else {
		return viper.Get(key)
	}
}

// 设置新配置值, 并将设置的值保存到配置文件
func SaveNewValue(key string, value any) {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		logger.Error("向配置文件保存设置时发生致命错误", "err", err.Error())
	}
}

// 设置默认的配置
func settingDefaultConfig() {
	// 自动启动应用时的默认设置
	viper.SetDefault("auto_startup.is_auto_run", false)
	viper.SetDefault("auto_startup.is_hide_windows", true)
	viper.SetDefault("auto_startup.is_hide_windows_old", true)

	// 手动打开应用时的默认设置
	viper.SetDefault("startup.is_hide_windows", false)
}

func createDefaultConfig() {
	settingDefaultConfig()
	if err := viper.SafeWriteConfig(); err != nil {
		logger.Error("创建默认配置文件时发生致命错误", "err", err.Error())
	}
}

// 将默认配置增量写入配置文件<不会影响配置文件中已有的配置>
func diffAndUpdateDefaultConfig() {
	settingDefaultConfig()
	if err := viper.WriteConfig(); err != nil {
		logger.Error("diff并增量更新默认配置至现有配置文件时发生致命错误", "err", err.Error())
	}
}
