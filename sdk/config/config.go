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

package config

import (
	"KeyTone/logger"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Store struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

var Clients_sse_stores sync.Map
var once_stores sync.Once

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

	// TIPS: 我们可能无法在此回调中, 直接获取被更改的配置是哪个
	//       > 因此如果使用此回调, 我们只需对我们所关注的配置项, 手动建立历史值, 并在回调中重新获取这个值予以对比即可。
	//       > (比如监听是否自动重启的配置是否更改, 如果更改就触发服务端推送<sse或webSocket>)
	viper.OnConfigChange(func(e fsnotify.Event) {
		// TIPS: 不要在这内部调用此函数或任何会修改配置文件的写入操作, 否则会造成循环依赖从而破坏事件
		// viper.WriteConfig()
		// TIPS: 虽然这里触发的sse在前端依旧会在前端触发写入操作(但前提是前端的值变了才会触发响应式的调用写入操作)
		//       * 因此, 我们在前端引入延时触发接收到sse后写入的逻辑(并做了进一步的防抖处理), 因此在sse进行写入时
		//         * 如果是由前端先触发写入, 并引起此处回调发送sse; 当前端延时一定时间后收到sse并作出同步时, 由于值已经变过了, 因此此sse的同步不会触发响应式更新, 也就打断了循环依赖。
		//           > 值得注意的是, 若是不进行延时, 由于网络拥堵因素, 可能造成旧值覆盖最新值的短时间循环依赖bug(直到碰到sse值和前端内存值一致后停止, 若是运气不好还有可能循环很多轮, 甚至无限循环)
		//         * 如果是有配置文件被手动变更引起的, 触发此处回调发送sse;  当前端延时一定时间后收到sse并作出同步时, 由于前端内存值还是旧值, 因此会被sse同步的最新值覆盖。(虽然会再次触发调用一次重复的写入, 并重新触发此处回调, 不过也就到此为止了, 循环会在前端再次收到后被打破。)
		go func(Clients_sse_stores *sync.Map) {
			stores := &Store{
				Key:   "get_all_value",
				Value: "",
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
	// * TIPS: viper在加载时, 会一次性将配置文件中所有配置读入内存中管理,也就是说,默认不会监听文件本身非viper操作之外的更改。
	//         > 若不调用WatchConfig(), 则只有通过viper的更改(如用`viper.Set`设置的更改), 才能够在下次`viper.Get`时读取到。
	viper.WatchConfig()

}

// 恢复默认设置(删除掉已有配置文件后, 再调用“createDefaultConfig()"函数就好了)
func RestoreDefaultSetting() {

}

var viperRWMutex sync.RWMutex

// 获取配置
// * 如果想获取这个json文件的所有配置, 则可用key="get_all_value"作为键值来查询
func GetValue(key string) any {
	viperRWMutex.RLock()
	defer viperRWMutex.RUnlock()
	if key == "get_all_value" {
		return viper.AllSettings()
	} else {
		return viper.Get(key)
	}
}

// 设置新配置值, 并将设置的值保存到配置文件
func SetValue(key string, value any) {
	viperRWMutex.Lock()
	defer viperRWMutex.Unlock()
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		logger.Error("向配置文件保存设置时发生致命错误", "err", err.Error())
	}

	// 由于viper.Set()在设计中拥有最高覆盖级别,因此需要在每次使用此api设置后, 清空viper.Set()的设置, 以使得文件监听的api可以正常工作。
	viper.Set(key, nil)

	// 等待 viper.WatchConfig 监听真实配置
	sleep := true
	ch := make(chan (struct{}))
	defer close(ch)

	go func(sleep *bool, ch chan struct{}) {
		// defer logger.Info("保护功能完成, 退出当前goroutine以结束保护--->应用配置项")
		for {
			select {
			case <-ch:
				// logger.Info("符合预期的退出行为--->应用配置项")
				return
			case <-time.After(time.Millisecond * 100): // 这个最大退出时间, 由您自由指定
				// logger.Warn("到达等待时间上限, 而进行的自动强制退出行为, 以避免资源浪费式的长期甚至永久等待行为--->应用配置项")
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
			// logger.Info("阻止了一次可能存在的错误删除行为--->应用配置项")
		}
	}
}

// 手动打开应用时的默认设置
const Startup___is_hide_windows = false

// 自动启动应用时的默认设置
const Auto_startup___is_auto_run = false
const Auto_startup___is_hide_windows = true
const Auto_startup___is_hide_windows_old = true

// 音频音量处理的默认设置
// * 用于设置页面 音量提升/缩减 设置
const Audio_volume_processing___volume_amplify = 0.0        // (-无穷)~0~(+无穷) ; // 此处理可能超出安全的范围, 它希望无论如何都要调整音量的大小(即使在一定值时会破坏原音频音质, 也在所不惜)
const Audio_volume_processing___volume_amplify_limit = 10.0 // 给(-无穷)~0~(+无穷)一个限制; 让其实际上为(-volume_amplify_limit)~0~(+volume_amplify_limit)// 设置默认的配置

// 主页面的默认设置
const Main_home___audio_volume_processing___volume_normal = 0.0                 // (-无穷)~0 ; // 让其最小值实际上为(-volume_amplify_limit)~0,以应对音量被增强后的音量缩减需求 // 此为安全范围内的音量处理, 它希望在保留(或不超过)原始音频音量的前提下调整音量
const Main_home___audio_volume_processing___volume_normal_reduce_scope = 5.0    // 默认为5.0
const Main_home___audio_volume_processing___volume_silent = false               // 当其为true时, 代表静音。
const Main_home___audio_volume_processing___is_open_volume_debug_slider = false // 用于在设置页面 和 主页面上显示 音量调试滑块

func settingDefaultConfig() {
	// 手动打开应用时的默认设置
	viper.SetDefault("startup.is_hide_windows", Startup___is_hide_windows)

	// 自动启动应用时的默认设置
	viper.SetDefault("auto_startup.is_auto_run", Auto_startup___is_auto_run)
	viper.SetDefault("auto_startup.is_hide_windows", Auto_startup___is_hide_windows)
	viper.SetDefault("auto_startup.is_hide_windows_old", Auto_startup___is_hide_windows_old)

	// 音频音量处理的默认设置
	// * 用于设置页面 音量提升/缩减 设置
	viper.SetDefault("audio_volume_processing.volume_amplify", Audio_volume_processing___volume_amplify)
	viper.SetDefault("audio_volume_processing.volume_amplify_limit", Audio_volume_processing___volume_amplify_limit)

	// 主页面的默认设置
	viper.SetDefault("main_home.audio_volume_processing.volume_normal", Main_home___audio_volume_processing___volume_normal)
	viper.SetDefault("main_home.audio_volume_processing.volume_normal_reduce_scope", Main_home___audio_volume_processing___volume_normal_reduce_scope)
	viper.SetDefault("main_home.audio_volume_processing.volume_silent", Main_home___audio_volume_processing___volume_silent)
	viper.SetDefault("main_home.audio_volume_processing.is_open_volume_debug_slider", Main_home___audio_volume_processing___is_open_volume_debug_slider)

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
