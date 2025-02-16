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

	// 设置配置文件名称和类型
	Viper.SetConfigName("config")
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
	v.SetConfigName("config")
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
