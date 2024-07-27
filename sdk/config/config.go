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
}

// 恢复默认设置
func RestoreDefaultSetting() {

}

// 获取配置
func GetValue(key string) any {
	// value := viper.GetString(key)
	value := viper.Get(key)
	return value
}

// 设置新配置值, 并将设置的值保存到配置文件
func SaveNewValue(key string, value any) {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		logger.Error("向配置文件保存设置时发生致命错误", "err", err.Error())
	}
}

// 创建默认设置的配置文件
func createDefaultSetting() {
	// 自动启动应用时的默认设置
	// 注意: 第二个参数不能是结构体, 而应是一个map[string]interface{}。
	// * 虽然interface{}可以表示map[string]interface{}, 但在go语言中, 直接使用interface{}来定义数据是编译器不可知不合法的。
	// viper.Set("autoStartup", struct {
	// 	isAutoRun     bool
	// 	isHideWindows bool
	// }{
	// 	isAutoRun:     false,
	// 	isHideWindows: true,
	// })
	// 第二个参数应该是一个map, 而不是一个结构体, 在go中我们使用map来映射json结构。
	// * 另外, 直接interface{} 定义键值对，这是不允许的, 会导致编译错误,因为 interface{} 本身不是一个可索引的类型。
	// * 不过, 反过来是可以使用 interface{}类型的变量, 来承载map[string]interface{}的, 只需在使用时进行"适当"的类型断言。
	//   * {
	//      	// 注意"适当", 当你不关注它时<如仅作为过程中的参数往后传递时>, 则无需断言它, 直接使用即可, 如案例中的变量aaa
	//      	var aaa any
	//      	bbb := map[string]interface{}{
	//      		"isAutoRun":     false,
	//      		"isHideWindows": true,
	//      	}
	//      	aaa = bbb
	//      	viper.Set("autoStartup", aaa)
	//
	//        // 因此, 下方内容是合法的
	//        {
	//     	   	 // 定义一个空的 map，用于存放解析后的数据, 是合法的
	//           // * 如果你不确定最终的json字符串的顶层结构是否对象(即映射)时
	//           var value any
	//           // 解析 JSON 字符串
	//           err = json.Unmarshal([]byte(store_setting.Value), &value)
	//        }
	//     }
	viper.Set("auto_startup", map[string]interface{}{
		"is_auto_run":         false,
		"is_hide_windows":     true,
		"is_hide_windows_old": true,
	})

	if err := viper.SafeWriteConfig(); err != nil {
		logger.Error("创建默认设置文件时发生致命错误", "err", err.Error())
	}
}

// 检查默认设置的键名是否都存在, 不存在则仅针对不存在的key, 创建对应默认配置(即不改变用户已有的自定义过的配置的前提下, 将新增的默认配置内容加入配置文件)。
func CheckDefaultSetting() {

}
