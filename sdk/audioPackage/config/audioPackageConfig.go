package audioPackageConfig

import (
	"KeyTone/logger"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 定义音频包根目录的路径
var AudioPackagePath string

var Viper *viper.Viper

// 定义Viper删除某个配置项时需要用到的全局变量
var deleteKeyValue []struct{} // 无法保证100%删除, 因为这种删除方式是利用了一些未知的黑盒bug。(tips: 千万不要完成其定义, 或者说千万不用对其进行任何形式的初始化。)

// 定义sse相关变量
type Store struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

var Clients_sse_stores sync.Map
var once_stores sync.Once

// 加载音频包时使用(也可用于创建新的音频包时)
func LoadConfig(configPath string, isCreate bool) {
	Viper = nil
	Viper = viper.New()

	// 设置配置文件名称和类型
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")

	// 添加配置文件路径
	Viper.AddConfigPath(configPath)

	// 读取配置文件
	if err := Viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			if isCreate {
				// 如果前端是选择创建的话, 属于已知预期行为,此时我们创建新的配置文件
				{
					// 创建新的配置文件前检查指定的路径是否存在
					if _, err := os.Stat(configPath); os.IsNotExist(err) {
						// 如果路径不存在，创建路径
						err := os.MkdirAll(configPath, os.ModePerm)
						if err != nil {
							logger.Error("配置文件路径创建时出错。", "err", err.Error())
						} else {
							logger.Info("配置文件路径创建成功。", "你的配置文件路径为", configPath)
						}
					} else if err != nil {
						logger.Error("检查配置文件路径时出错。", "err", err.Error())
					} else {
						logger.Info("配置文件路径已存在且无异常。", "你的配置文件路径为", configPath)
					}
				}
				createDefaultConfig()
			} else {
				// 否则为正常的加载, 默认不会创建配置文件, 以识别出有问题的键音包
				logger.Error("未找到正确的音频包配置", "path", configPath)
				// TODO: 可以返回给前端, 供其提示用户
				// SSE
			}
		} else {
			// 其他错误
			logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
			// TODO: 可以返回给前端, 供其提示用户
			// SSE
		}
	} else {
		logger.Info("音频包已加载, 正在与DefaultConfig进行diff和增量载入...")
		// 如果正常加载了键音配置文件, 则进行增量式的检测与更新(以在键音包设置出现更新时, 最大程度的兼容旧版本)
		diffAndUpdateDefaultConfig()
		logger.Info("音频包diff和增量载入完成")
	}

	Viper.OnConfigChange(func(e fsnotify.Event) {
		go func(Clients_sse_stores *sync.Map) {
			stores := &Store{
				Key:   "get_all_value",
				Value: GetValue("get_all_value"),
			}
			Clients_sse_stores.Range(func(key, value interface{}) bool {
				clientChan := key.(chan *Store)
				serverChan := value.(chan bool)
				select {
				case clientChan <- stores:
					return true
				case <-serverChan:
					once_stores.Do(func() {
						close(serverChan)
					})
					return true
				}
			})
		}(&Clients_sse_stores)
	})

	// 监听配置文件更改
	Viper.WatchConfig()

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

	// 等待 viper.WatchConfig 监听真实配置
	sleep := true
	ch := make(chan (struct{}))
	defer close(ch)

	go func(sleep *bool, ch chan struct{}) {
		defer logger.Info("保护功能完成, 退出当前goroutine以结束保护--->音频包配置项")
		for {
			select {
			case <-ch:
				logger.Info("符合预期的退出行为--->音频包配置项")
				return
			case <-time.After(time.Millisecond * 100): // 这个最大退出时间, 由您自由指定
				logger.Warn("到达等待时间上限, 而进行的自动强制退出行为, 以避免资源浪费式的长期甚至永久等待行为--->音频包配置项")
				*sleep = false
				return
			}
		}
	}(&sleep, ch)

	for sleep {
		if viper.Get(key) != nil {
			sleep = false
			// 如果函数结束, 通道自然会关闭, 从而解除阻塞行为。无需使用下行中可能引入新阻塞的逻辑。
			// ch <- struct{}{}
		} else {
			logger.Info("阻止了一次可能存在的错误删除行为--->音频包配置项")
		}
	}
}

// 删除某个配置项
func DeleteValue(key string) {
	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)

	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)

	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)

	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)
}

/**********************************************************************************************************************/

// var CreateViper *viper.Viper // 创建过程中, 播放的音频包理应是当前正在制作的。

// 音频包的名字, 不应该由目录名决定, 或者说目录名甚至可以是任意随机的UUID。
// * 目录名的映射由前端处理, 前端将会在加载用户历史选择音频包前, 先遍历所有音频包, 以建立音频包名称与路径的映射。
//   > 比如变量过程中, 可以反复使用LoadConfig和GetValue来获取每个 音频包的名称。
// * 不允许没有名字的音频包出现, 前端在创建映射以生成列表的过程中, 默认忽略没有名字的音频包。(即mvvm相关的map结构体/对象中, 不应有这些)
//   * 当然, 在前端很好避免, 加一个输入窗口校验即可。忽略策略是应对极少情况下用户更改配置文件删掉了音频包名字的情况。
//   * 因此, 会给一个默认名称, 比如: 新的音频包(New audio package)。

// 创建过程中, 是默认保存的, 无需手动保存。 如果音频包不想要了, 需要手动删除。

// 键音包名称
const Package_name = "新的键音包" // 其实这个在此处无关紧要, 因为此默认名称的设置, 主要还是在前端进行, 以前端国际化为主。

func settingDefaultConfig() {
	// 手动打开应用时的默认设置
	Viper.SetDefault("package_name", Package_name)

}

func createDefaultConfig() {
	settingDefaultConfig()
	if err := Viper.SafeWriteConfig(); err != nil {
		logger.Error("创建默认音频包配置文件时发生致命错误", "err", err.Error())
	}
}

// 将默认配置增量写入配置文件<不会影响配置文件中已有的配置>
func diffAndUpdateDefaultConfig() {
	settingDefaultConfig()
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("diff并增量更新默认音频包配置至现有音频包配置文件时发生致命错误", "err", err.Error())
	}
}
