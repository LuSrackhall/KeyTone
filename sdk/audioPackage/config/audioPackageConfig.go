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

package audioPackageConfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"KeyTone/audioPackage/enc"
	"KeyTone/logger"

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

var viperRWMutex sync.RWMutex

const (
	encryptionModePlain     = "plain"
	encryptionModeLegacyHex = "legacy-hex"
	encryptionModeCore      = "core"
)

type encryptionState struct {
	Mode      string
	AlbumPath string
	CorePath  string
	StubPath  string
}

var currentEncState encryptionState

func resetEncryptionState() {
	currentEncState = encryptionState{}
}

func setEncryptionState(mode, albumPath, corePath, stubPath string) {
	currentEncState = encryptionState{
		Mode:      mode,
		AlbumPath: albumPath,
		CorePath:  corePath,
		StubPath:  stubPath,
	}
}

func loadConfigFromCore(configPath string, stub *coreStubMetadata) error {
	albumUUID := filepath.Base(configPath)
	coreFile := filepath.Join(configPath, stub.Core)

	cipherBytes, err := os.ReadFile(coreFile)
	if err != nil {
		return fmt.Errorf("读取 core 文件失败: %w", err)
	}

	plainJSON, err := enc.DecryptConfigBytes(cipherBytes, albumUUID)
	if err != nil {
		return fmt.Errorf("解密 core 文件失败: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "keytone-album-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}

	tmpAlbumDir := filepath.Join(tmpDir, albumUUID)
	if err := os.MkdirAll(tmpAlbumDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建临时专辑目录失败: %w", err)
	}

	tmpPkg := filepath.Join(tmpAlbumDir, "package.json")
	if err := os.WriteFile(tmpPkg, []byte(plainJSON), 0644); err != nil {
		return fmt.Errorf("写入临时配置失败: %w", err)
	}

	Viper = viper.New()
	Viper.SetConfigName("package")
	Viper.SetConfigType("json")
	Viper.AddConfigPath(tmpAlbumDir)
	if err := Viper.ReadInConfig(); err != nil {
		return fmt.Errorf("加载临时配置失败: %w", err)
	}

	logger.Info("音频包通过 core 文件解密加载成功", "album", albumUUID)
	diffAndUpdateDefaultConfig()
	setEncryptionState(encryptionModeCore, configPath, coreFile, filepath.Join(configPath, "package.json"))
	// 注意：attachConfigWatcher() 由调用者在锁释放后执行，避免死锁
	return nil
}

func attachConfigWatcher() {
	Viper.OnConfigChange(func(e fsnotify.Event) {
		go func(Clients_sse_stores *sync.Map) {
			viperRWMutex.RLock()
			if Viper == nil {
				viperRWMutex.RUnlock()
				return
			}
			stores := &Store{
				Key:   "get_all_value",
				Value: Viper.AllSettings(),
			}
			viperRWMutex.RUnlock()
			if stores.Value == nil {
				return
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

	Viper.WatchConfig()
}

// migrateConfigFile 自动检测并迁移旧的 config.json 为 package.json
// 用于支持从旧版本到新版本的平滑升级
func migrateConfigFile(configPath string) error {
	oldConfigPath := filepath.Join(configPath, "config.json")
	newConfigPath := filepath.Join(configPath, "package.json")

	// 检查新的配置文件是否已存在
	if _, err := os.Stat(newConfigPath); err == nil {
		// 新文件已存在，检查旧文件是否也存在
		if _, err := os.Stat(oldConfigPath); err == nil {
			// 新旧文件都存在，这是一个异常情况，记录警告但不处理
			logger.Warn("检测到 config.json 和 package.json 同时存在，将保留 package.json", "path", configPath)
			return nil
		}
		// 新文件存在，旧文件不存在，无需迁移
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误（如权限问题）
		logger.Error("检查新配置文件时出错", "path", newConfigPath, "error", err.Error())
		return err
	}

	// 检查旧的配置文件是否存在
	if _, err := os.Stat(oldConfigPath); err == nil {
		// 旧文件存在，执行迁移
		if err := os.Rename(oldConfigPath, newConfigPath); err != nil {
			logger.Error("迁移配置文件失败", "from", oldConfigPath, "to", newConfigPath, "error", err.Error())
			return err
		}
		logger.Info("配置文件成功迁移", "from", "config.json", "to", "package.json", "path", configPath)
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误
		logger.Error("检查旧配置文件时出错", "path", oldConfigPath, "error", err.Error())
		return err
	}

	// 新旧文件都不存在，无需迁移
	return nil
}

func LoadConfig(configPath string, isCreate bool) {
	// Viper重新初始化的过程, 是属于临界区的, 因此需要加锁。
	// 注意：attachConfigWatcher() 必须在锁释放后调用，以避免潜在死锁。
	// 使用 shouldAttachWatcher 标记来延迟 watcher 的附加。
	var shouldAttachWatcher bool
	defer func() {
		if shouldAttachWatcher {
			attachConfigWatcher()
		}
	}()

	viperRWMutex.Lock()
	defer viperRWMutex.Unlock()
	if Viper != nil {
		Viper.StopWatch()
		Viper = nil
	}
	Viper = viper.New()
	resetEncryptionState()

	// 尝试迁移旧的 config.json 文件为 package.json
	if err := migrateConfigFile(configPath); err != nil {
		logger.Warn("配置文件迁移检查失败，继续使用现有配置", "error", err.Error())
	}

	pkgPath := filepath.Join(configPath, "package.json")
	stubMeta, pkgRaw, stubErr := readCoreStub(configPath)
	if stubErr != nil {
		logger.Error("解析加密指示 JSON 失败", "err", stubErr.Error())
		Viper = nil
		return
	}
	if stubMeta != nil {
		if err := loadConfigFromCore(configPath, stubMeta); err != nil {
			logger.Error("core 文件加载失败", "err", err.Error())
			Viper = nil
		} else {
			shouldAttachWatcher = true
		}
		return
	}

	// 设置配置文件名称和类型（使用新的 package.json）
	Viper.SetConfigName("package")
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
				// 如果未找到正确的音频包配置, 则应该清空Viper, 防止内存中残留的Viper被错误的使用(比如前端读取配置, 但实际上配置文件并不存在, 却读取到了内存中的残留配置)
				Viper = nil
				// TODO: 可以返回给前端, 供其提示用户
				// SSE
			}
		} else {
			// 其他错误（包括 JSON 解析失败），尝试解密回退
			logger.Warn("读取音频包配置失败，尝试解密回退", "err", err.Error())
			// 尝试读取 package.json 原始内容
			b := pkgRaw
			var readErr error
			if b == nil {
				b, readErr = os.ReadFile(pkgPath)
			}
			if readErr != nil {
				logger.Error("读取package.json失败，无法解密回退", "err", readErr.Error())
				Viper = nil
				return
			}
			// 通过目录名作为 albumUUID（默认约定：目录名即专辑UUID），若配置中已有字段则导入后覆盖
			albumUUID := filepath.Base(configPath)
			if enc.IsLikelyHexCipher(b) {
				plain, decErr := enc.DecryptConfigHex(strings.TrimSpace(string(b)), albumUUID)
				if decErr != nil {
					logger.Error("解密回退失败", "err", decErr.Error())
					Viper = nil
					return
				}
				// 将解密后的JSON写入临时目录供Viper加载
				tmpDir, terr := os.MkdirTemp("", "keytone-album-*")
				if terr != nil {
					logger.Error("创建临时目录失败", "err", terr.Error())
					Viper = nil
					return
				}
				tmpAlbumDir := filepath.Join(tmpDir, albumUUID)
				_ = os.MkdirAll(tmpAlbumDir, os.ModePerm)
				tmpPkg := filepath.Join(tmpAlbumDir, "package.json")
				if werr := os.WriteFile(tmpPkg, []byte(plain), 0644); werr != nil {
					logger.Error("写入临时配置失败", "err", werr.Error())
					Viper = nil
					return
				}
				// 重新指向临时目录加载
				Viper = viper.New()
				Viper.SetConfigName("package")
				Viper.SetConfigType("json")
				Viper.AddConfigPath(tmpAlbumDir)
				if r2 := Viper.ReadInConfig(); r2 != nil {
					logger.Error("临时配置加载失败", "err", r2.Error())
					Viper = nil
					return
				}
				logger.Info("音频包通过解密回退加载成功", "album", albumUUID)
				// 增量写默认项
				diffAndUpdateDefaultConfig()
				setEncryptionState(encryptionModeLegacyHex, configPath, "", pkgPath)
				shouldAttachWatcher = true
				// 在每次明确写入后执行回写加密（在SetValue完成后触发已存在），此处确保首次加载时也能写回一次（不强制）
				// 不主动覆写源文件以避免无改动写入；实际回写在 SetValue 调用中实现
				return
			}
			// 非十六进制密文但解析失败，按原逻辑报错
			logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
		}
	} else {
		logger.Info("音频包已加载, 正在与DefaultConfig进行diff和增量载入...")
		// 如果正常加载了键音配置文件, 则进行增量式的检测与更新(以在键音包设置出现更新时, 最大程度的兼容旧版本)
		diffAndUpdateDefaultConfig()
		logger.Info("音频包diff和增量载入完成")
		setEncryptionState(encryptionModePlain, configPath, "", pkgPath)
	}

	if Viper != nil {
		if currentEncState.Mode == "" {
			setEncryptionState(encryptionModePlain, configPath, "", pkgPath)
		}
		shouldAttachWatcher = true
	}

}

func GetValue(key string) any {
	viperRWMutex.RLock()
	defer viperRWMutex.RUnlock()
	if Viper == nil {
		return nil
	}
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
	if Viper == nil {
		return
	}
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
		// defer logger.Info("保护功能完成, 退出当前goroutine以结束保护--->音频包配置项")
		for {
			select {
			case <-ch:
				// logger.Info("符合预期的退出行为--->音频包配置项")
				return
			case <-time.After(time.Millisecond * 100): // 这个最大退出时间, 由您自由指定
				// logger.Warn("到达等待时间上限, 而进行的自动强制退出行为, 以避免资源浪费式的长期甚至永久等待行为--->音频包配置项")
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
			// logger.Info("阻止了一次可能存在的错误删除行为--->音频包配置项")
		}
	}

	// 若源文件为加密形态，则需要将临时目录的明文重新加密写回源文件
	albumPath := currentEncState.AlbumPath
	if albumPath == "" {
		albumPath = AudioPackagePath
	}
	pkgPath := filepath.Join(albumPath, "package.json")
	mode := currentEncState.Mode

	switch mode {
	case encryptionModeCore:
		tmpPlain := Viper.ConfigFileUsed()
		plainBytes, err := os.ReadFile(tmpPlain)
		if err != nil {
			logger.Error("读取临时配置失败", "err", err.Error())
			return
		}
		albumUUID := filepath.Base(albumPath)
		cipherBytes, err := enc.EncryptConfigBytes(string(plainBytes), albumUUID)
		if err != nil {
			logger.Error("加密 core 配置失败", "err", err.Error())
			return
		}
		corePath := currentEncState.CorePath
		if corePath == "" {
			corePath = filepath.Join(albumPath, CoreFileName)
		}
		tmpOut := corePath + ".tmp"
		if err := os.WriteFile(tmpOut, cipherBytes, 0644); err != nil {
			logger.Error("写入 core 临时文件失败", "err", err.Error())
			return
		}
		if err := os.Rename(tmpOut, corePath); err != nil {
			_ = os.Remove(tmpOut)
			logger.Error("重命名 core 临时文件失败", "err", err.Error())
			return
		}
		if err := writeCoreStub(albumPath, nil); err != nil {
			logger.Error("更新指示 JSON 失败", "err", err.Error())
		}

	case encryptionModeLegacyHex:
		albumUUID := filepath.Base(albumPath)
		tmp := Viper.ConfigFileUsed()
		b, rerr := os.ReadFile(tmp)
		if rerr != nil {
			logger.Error("读取临时配置失败", "err", rerr.Error())
			return
		}
		cipherHex, eerr := enc.EncryptConfigJSON(string(b), albumUUID)
		if eerr != nil {
			logger.Error("加密配置写回失败", "err", eerr.Error())
			return
		}
		tmpOut := pkgPath + ".tmp"
		if werr := os.WriteFile(tmpOut, []byte(cipherHex), 0644); werr != nil {
			logger.Error("写入加密回写临时文件失败", "err", werr.Error())
			return
		}
		if err := os.Rename(tmpOut, pkgPath); err != nil {
			_ = os.Remove(tmpOut)
			logger.Error("重命名 package 临时文件失败", "err", err.Error())
		}

	default:
		// 兼容旧逻辑：若仍检测到十六进制密文则执行旧写回
		if f, err := os.Open(pkgPath); err == nil {
			defer f.Close()
			buf, _ := io.ReadAll(f)
			if enc.IsLikelyHexCipher(buf) {
				albumUUID := filepath.Base(albumPath)
				tmp := Viper.ConfigFileUsed()
				b, rerr := os.ReadFile(tmp)
				if rerr == nil {
					cipherHex, eerr := enc.EncryptConfigJSON(string(b), albumUUID)
					if eerr == nil {
						tmpOut := pkgPath + ".tmp"
						if werr := os.WriteFile(tmpOut, []byte(cipherHex), 0644); werr == nil {
							_ = os.Rename(tmpOut, pkgPath)
						} else {
							logger.Error("写入加密回写临时文件失败", "err", werr.Error())
						}
					} else {
						logger.Error("加密配置写回失败", "err", eerr.Error())
					}
				}
			}
		}
	}
}

// 删除某个配置项
func DeleteValue(key string) {
	viperRWMutex.Lock()
	defer viperRWMutex.Unlock()
	if Viper == nil {
		return
	}
	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)
	time.Sleep(time.Millisecond * 10)

	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)
	time.Sleep(time.Millisecond * 10)

	Viper.Set(key, deleteKeyValue)
	if err := Viper.WriteConfig(); err != nil {
		logger.Error("删除键值对时发生致命错误", "err", err.Error())
	}
	Viper.Set(key, nil)
	time.Sleep(time.Millisecond * 10)

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

// 这里没必要判断 Viper == nil, 因为这个函数是在创建新音频包时调用的, 一定是Viper != nil的情况。(同理, 这里也没必要加锁。)
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

// ClearConfig 清除当前的 Viper 配置实例
// 这是一个线程安全的操作，会获取写锁
func ClearConfig() {
	viperRWMutex.Lock()
	defer viperRWMutex.Unlock()
	if Viper != nil {
		Viper = nil
	}
}
