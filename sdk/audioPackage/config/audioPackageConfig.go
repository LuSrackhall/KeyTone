package audioPackageConfig

import (
	"KeyTone/logger"
	"sync"

	"github.com/spf13/viper"
)

// 定义音频包根目录的路径
var AudioPackagePath string

var Viper *viper.Viper

// 加载现有音频包时使用
func LoadConfig(configPath string) {
	Viper = nil
	Viper = viper.New()

	// 设置配置文件名称和类型
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")

	// 添加配置文件路径
	Viper.AddConfigPath(configPath)

	// 监听配置文件更改
	Viper.WatchConfig()

	// 读取配置文件
	if err := Viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			logger.Error("未找到正确的音频包配置", "path", configPath)
			// TODO: 可以返回给前端, 供其提示用户
			// SSE
		} else {
			// 其他错误
			logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
			// TODO: 可以返回给前端, 供其提示用户
			// SSE
		}
	}
}

var viperRWMutex sync.RWMutex

func GetValue(key string) any {
	viperRWMutex.RLock()
	defer viperRWMutex.RUnlock()
	if key == "get_all_value" {
		return Viper.AllSettings()
	} else {
		return Viper.Get(key)
	}
}

// 设置新配置值, 并将设置的值保存到配置文件
func SetValue(key string, value any) {
	viperRWMutex.Lock()
	defer viperRWMutex.Unlock()
	Viper.Set(key, value)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("向音频包保存配置时发生致命错误", "err", err.Error())
	}

	// 由于viper.Set()在设计中拥有最高覆盖级别,因此需要在每次使用此api设置后, 清空viper.Set()的设置, 以使得文件监听的api可以正常工作。
	Viper.Set(key, nil)
}

/**********************************************************************************************************************/

// var CreateViper *viper.Viper // 创建过程中, 播放的音频包理应是当前正在制作的。

// 创建新的音频包时使用
func CreateConfig(configPath string) {}

// 音频包的名字, 不应该由目录名决定, 或者说目录名甚至可以是任意随机的UUID。
// * 目录名的映射由前端处理, 前端将会在加载用户历史选择音频包前, 先遍历所有音频包, 以建立音频包名称与路径的映射。
//   > 比如变量过程中, 可以反复使用LoadConfig和GetValue来获取每个 音频包的名称。
// * 不允许没有名字的音频包出现, 前端在创建映射以生成列表的过程中, 默认忽略没有名字的音频包。(即mvvm相关的map结构体/对象中, 不应有这些)
//   * 当然, 在前端很好避免, 加一个输入窗口校验即可。忽略策略是应对极少情况下用户更改配置文件删掉了音频包名字的情况。
//   * 因此, 会给一个默认名称, 比如: 新的音频包(New audio package)。

// 创建过程中, 是默认保存的, 无需手动保存。 如果音频包不想要了, 需要手动删除。
