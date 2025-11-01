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

package audioPackageList

import (
	"KeyTone/logger"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

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

// 用于获取键音包根目录中, 所有键音包的完整路径的列表。
func GetAudioPackageList(rootDir string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(rootDir)

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(rootDir, entry.Name()))
		}
	}

	return directories, err
}

// 用于获取具体路径下的键音包的名称。
var viperRWMutex sync.RWMutex

func GetAudioPackageName(configPath string) any {
	// 加载配置文件
	Viper := viper.New()
	defer func() {
		Viper = nil
	}()

	// 尝试迁移旧的 config.json 文件为 package.json
	if err := migrateConfigFile(configPath); err != nil {
		logger.Warn("配置文件迁移检查失败，继续使用现有配置", "error", err.Error())
	}

	// 设置配置文件名称和类型
	Viper.SetConfigName("package")
	Viper.SetConfigType("json")

	// 添加配置文件路径
	Viper.AddConfigPath(configPath)

	// 读取配置文件
	if err := Viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 否则为正常的加载, 默认不会创建配置文件, 以识别出有问题的键音包
			logger.Error("未找到正确的音频包配置", "path", configPath)
		} else {
			// 其他错误
			logger.Error("读取音频包配置时发生致命错误", "err", err.Error())
		}
	}

	// 从加载的配置文件中, 读取 package_name 字段
	viperRWMutex.RLock()
	defer viperRWMutex.RUnlock()
	return Viper.Get("package_name")
}

// 更新专辑配置文件中的 UUID
func UpdateAlbumUUID(albumPath string, newUUID string) error {
	// 加载配置文件
	v := viper.New()
	v.SetConfigName("package")
	v.SetConfigType("json")
	v.AddConfigPath(albumPath)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 更新 UUID
	v.Set("audio_pkg_uuid", newUUID)

	// 保存更改
	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}
