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

package server

import (
	audioPackageConfig "KeyTone/audioPackage/config"
	audioPackageList "KeyTone/audioPackage/list"
	"KeyTone/config"
	"KeyTone/keyEvent"
	"KeyTone/keySound"
	"KeyTone/logger"
	"archive/zip"
	"bytes"
	"crypto"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type KeyStateMessage struct {
	Type    string `json:"type"`
	Keycode uint16 `json:"keycode"`
	State   string `json:"state"`
}

func ServerRun() {
	// 启动gin
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	mainRouters(r)

	r.GET("/stream", func(c *gin.Context) {

		logger.Logger.Debug("新生成了一个线程............................")
		logger.Debug("新生成了一个线程............................")

		clientStoresChan := make(chan *config.Store)
		clientAudioPackageStoresChan := make(chan *audioPackageConfig.Store)
		clientKeyEventStoresChan := make(chan *keyEvent.Store)

		serverStoresChan := make(chan bool, 1)
		serverAudioPackageStoresChan := make(chan bool, 1)
		serverKeyEventStoresChan := make(chan bool, 1)

		config.Clients_sse_stores.Store(clientStoresChan, serverStoresChan)
		audioPackageConfig.Clients_sse_stores.Store(clientAudioPackageStoresChan, serverAudioPackageStoresChan)
		keyEvent.Clients_sse_stores.Store(clientKeyEventStoresChan, serverKeyEventStoresChan)

		defer func() {
			config.Clients_sse_stores.Delete(clientStoresChan)
			audioPackageConfig.Clients_sse_stores.Delete(clientAudioPackageStoresChan)
			keyEvent.Clients_sse_stores.Delete(clientKeyEventStoresChan)

			logger.Logger.Debug("一个线程退出了............................")
			logger.Debug("一个线程退出了............................")
		}()

		clientGone := c.Request.Context().Done()

		for {
			re := c.Stream(func(w io.Writer) bool {
				select {
				case <-clientGone:
					serverStoresChan <- false
					serverAudioPackageStoresChan <- false
					serverKeyEventStoresChan <- false

					return false

				case message, ok := <-clientStoresChan:
					if !ok {
						logger.Error("通道clientStoresChan非正常关闭")
						return true
					}
					c.SSEvent("message", message)
					return true
				case messageAudioPackage, ok := <-clientAudioPackageStoresChan:
					if !ok {
						logger.Error("通道clientAudioPackageStoresChan非正常关闭")
						return true
					}
					c.SSEvent("messageAudioPackage", messageAudioPackage)
					return true
				case messageKeyEvent, ok := <-clientKeyEventStoresChan:
					if !ok {
						logger.Error("通道clientKeyEventStoresChan非正常关闭")
						return true
					}
					c.SSEvent("messageKeyEvent", messageKeyEvent)
					return true
				}
			})

			if !re {
				return
			}
		}

	})

	keytonePkgRouters(r)

	// 尝试在指定端口启动服务
	listener, err := net.Listen("tcp", "localhost:38888")
	if err != nil {
		// 如果38888被占用，让系统分配一个可用端口
		listener, err = net.Listen("tcp", "localhost:0")
		if err != nil {
			logger.Error("无法启动服务:", err)
			return
		}
	}

	// 获取实际使用的端口
	port := listener.Addr().(*net.TCPAddr).Port

	// 输出端口信息，让Electron主进程可以捕获
	fmt.Printf("KEYTONE_PORT=%d\n", port)

	// 使用listener启动服务
	go r.RunListener(listener)
}

func mainRouters(r *gin.Engine) {
	settingStoreRouters := r.Group("/store")

	// 给到'客户端'或'前端'使用, 供它们获取持久化的设置。
	settingStoreRouters.GET("/get", func(ctx *gin.Context) {

		// key := ctx.Query("key")
		key := ctx.DefaultQuery("key", "unknown")

		if key == "unknown" || key == "" {
			ctx.JSON(200, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:",
			})
			return
		}

		value := config.GetValue(key)

		fmt.Println("查询到的value= ", value)

		ctx.JSON(200, gin.H{
			"message": "ok",
			"key":     key,
			// 这里的value, 会自动转换为JSON字符串
			"value": value,
		})
	})

	// 给到'前端'使用, 供其ui界面实现应用的设置功能
	settingStoreRouters.POST("/set", func(ctx *gin.Context) {
		type SettingStore struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		}

		var store_setting SettingStore
		err := ctx.ShouldBind(&store_setting)
		if err != nil || store_setting.Key == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		config.SetValue(store_setting.Key, store_setting.Value)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

}

func keytonePkgRouters(r *gin.Engine) {

	keytonePkgRouters := r.Group("/keytone_pkg")

	// 加载键音包
	keytonePkgRouters.POST("/load_config", func(ctx *gin.Context) {
		type Arg struct {
			AudioPkgUUID string `json:"audioPkgUUID"`
			IsCreate     bool   `json:"isCreate"`
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		var audioPkgPath string
		if arg.IsCreate {
			// 如果创建新的键音包, 我们只知道最终的uuid(不知道存放uuid文件夹的路径)
			audioPkgPath = filepath.Join(audioPackageConfig.AudioPackagePath, arg.AudioPkgUUID)
		} else {
			// 如果加载已有的键音包, 我们知道键音包的完整路径(包括uuid与即uuid文件夹所在的路径)
			audioPkgPath = arg.AudioPkgUUID
		}

		// 加载键音包配置文件
		audioPackageConfig.LoadConfig(audioPkgPath, arg.IsCreate)

		ctx.JSON(200, gin.H{
			"message":      "ok",
			"audioPkgPath": audioPkgPath,
		})
	})

	// 接收前端上传的音频文件, 并存入本地路径
	keytonePkgRouters.POST("/add_new_sound_file", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件添加失败, 传输问题:" + err.Error(),
			})
			return
		}

		// 打开上传的文件
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 无法打开上传的文件:" + err.Error(),
			})
			return
		}
		defer src.Close()

		// 计算文件的SHA256哈希值
		hash := crypto.SHA256.New()
		if _, err := io.Copy(hash, src); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 计算文件哈希值失败:" + err.Error(),
			})
			return
		}
		hashSum := hash.Sum(nil)
		hashString := fmt.Sprintf("%x", hashSum)

		// 获取文件扩展名
		ext := filepath.Ext(file.Filename)

		// 使用哈希值作为文件名
		newFileName := hashString + ext

		// 测试此文件是否已经存在(有些相同的文件, 可能文件名称不同。) 若不存在则获取结果为nil, 可以正常往下走。 若存在则获取结果不为nil, 仅需向后添加文件名即可返回给前端, 无需后续步骤中重复保存相同的文件。(但对于前端用户, 不影响其认为这是两个不同的文件, 因为我们对于名称单独进行了保存)
		// * 至于文件名称重复的问题, 此处不作处理, 皆由用户自行管理名称, 不只是同一sha256uuid的名字可可重复, 甚至允许用户对不同sha256uuid的音频文件起相同的名字, 皆由用户自由发挥即可。
		if audioPackageConfig.GetValue("audio_files."+hashString) != nil {
			count := 0
			for audioPackageConfig.GetValue("audio_files."+hashString+".name."+strconv.Itoa(count)) != nil {
				count++
			}
			audioPackageConfig.SetValue("audio_files."+hashString, map[string]any{
				/**
				 * filepath.Base(file.Filename)：
				 *	- 这个函数返回路径中的最后一个元素（文件名）。
				 *	- 例如，如果 file.Filename 是 "/path/to/myFile.txt"，这个函数会返回 "myFile.txt"。
				 *	filepath.Ext(file.Filename)：
				 *	- 这个函数返回文件名的扩展名，包括点号。
				 *	- 对于 "myFile.txt"，它会返回 ".txt"。
				 *	strings.TrimSuffix(base, ext)：
				 *	- 这个函数从第一个参数（base）的末尾移除第二个参数（ext）指定的后缀。
				 *	- 如果 base 是 "myFile.txt"，ext 是 ".txt"，结果就是 "myFile"。
				 */
				"name": map[string]any{
					strconv.Itoa(count): strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename)), // strings.Split(file.Filename, ".")[0]
				},
				"type": ext,
			})

			// 因文件已存在与文件系统中, 故无需继续进行真实的文件保存。 这里直接将正确完成的消息返回给前端, 并退出此次请求的处理即可。
			ctx.JSON(200, gin.H{
				"message":  "ok",
				"fileName": newFileName,
			})

			// 退出此次请求的处理 (TIPS: 单纯的向前端返回消息, 并不能自动return。 此处我们需要主动退出, 防止执行后续步骤造成画蛇添足。)
			return
		}

		audioPkgUUID, ok := audioPackageConfig.GetValue("audio_pkg_uuid").(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取音频包UUID失败",
			})
			return
		}

		// 保存文件
		destPath := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", newFileName)
		if err := ctx.SaveUploadedFile(file, destPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 文件添加失败, 后端保存过程中发生错误:" + err.Error(),
			})
			return
		}

		// 文件如果可以成功保存, 证明这是首次未重复过的文件, 因此 使用"0"作为名字的key值。(这是为了应对用户上传相同音频文件时但不想无辜增大键音包的策略)
		// 文件保存成功后, 将原文件名作为value值(裁掉扩展名,只要文件名字), sha256哈希值文件名作为key值(裁掉扩展名), 存入键音包配置文件中的audioFiles对象中。
		// 源文件名作为value值, 是因为key值中不允许大写字符出现, 因此不能应对用户对音频名称的复杂设置需求。而且, 它本身也应该是作为value值存储的。
		// 哈希值作为key值, 也刚好符合sha256哈希值通常用纯小写表示的惯例。至于真实文件后缀或者说文件类型, 则也存储至value中去。
		// audioPackageConfig.SetValue("audio_files."+hashString+".name", strings.Split(file.Filename, ".")[0])
		// audioPackageConfig.SetValue("audio_files."+hashString+".type", ext)
		audioPackageConfig.SetValue("audio_files."+hashString, map[string]any{
			/**
			 * filepath.Base(file.Filename)：
			 *	- 这个函数返回路径中的最后一个元素（文件名）。
			 *	- 例如，如果 file.Filename 是 "/path/to/myFile.txt"，这个函数会返回 "myFile.txt"。
			 *	filepath.Ext(file.Filename)：
			 *	- 这个函数返回文件名的扩展名，包括点号。
			 *	- 对于 "myFile.txt"，它会返回 ".txt"。
			 *	strings.TrimSuffix(base, ext)：
			 *	- 这个函数从第一个参数（base）的末尾移除第二个参数（ext）指定的后缀。
			 *	- 如果 base 是 "myFile.txt"，ext 是 ".txt"，结果就是 "myFile"。
			 */
			"name": map[string]any{
				"0": strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename)), // strings.Split(file.Filename, ".")[0]
			},
			"type": ext,
		})

		// 全部处理完毕后, 将正确完成的消息返回给前端
		ctx.JSON(200, gin.H{
			"message":  "ok",
			"fileName": newFileName,
		})
	})

	keytonePkgRouters.GET("/get", func(ctx *gin.Context) {

		// key := ctx.Query("key")
		key := ctx.DefaultQuery("key", "unknown")

		if key == "unknown" || key == "" {
			ctx.JSON(200, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:",
			})
			return
		}

		value := audioPackageConfig.GetValue(key)

		fmt.Println("查询到的value= ", value)

		ctx.JSON(200, gin.H{
			"message": "ok",
			"key":     key,
			// 这里的value, 会自动转换为JSON字符串
			"value": value,
		})
	})

	keytonePkgRouters.POST("/set", func(ctx *gin.Context) {
		type SettingStore struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		}

		var store_setting SettingStore
		err := ctx.ShouldBind(&store_setting)
		if err != nil || store_setting.Key == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		audioPackageConfig.SetValue(store_setting.Key, store_setting.Value)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	keytonePkgRouters.POST("/delete", func(ctx *gin.Context) {
		// 前端使用 axios 发送 POST 请求时使用了 JSON 格式而不是 form 格式,
		// 所以需要使用 ShouldBind 来绑定 JSON 数据而不是 PostForm
		type DeleteArg struct {
			Key string `json:"key"`
		}

		var arg DeleteArg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.Key == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容key值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		audioPackageConfig.DeleteValue(arg.Key)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	keytonePkgRouters.POST("/sound_file_rename", func(ctx *gin.Context) {
		type Arg struct {
			Sha256 string `json:"sha256"`
			NameID string `json:"nameID"`
			Name   string `json:"name"`
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.Sha256 == "" || arg.NameID == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		audioPackageConfig.SetValue("audio_files."+arg.Sha256+".name."+arg.NameID, arg.Name)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	keytonePkgRouters.POST("/sound_file_delete", func(ctx *gin.Context) {
		type Arg struct {
			Sha256 string `json:"sha256"` // 文件名ID(实际文件名)
			NameID string `json:"nameID"` // 文件名ID(UI端使用, 用于索引虚拟文件名)
			Type   string `json:"type"`   // 文件类型
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.Sha256 == "" || arg.NameID == "" || arg.Type == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			logger.Error("message", "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:"+err.Error())
			return
		}

		audioPackageConfig.DeleteValue("audio_files." + arg.Sha256 + ".name." + arg.NameID)

		// TIPS: 每次删除操作后, 都清除内存中的name字段, 并依赖viper提供的实时更新特性与实际文件保持一致。
		// 			 * 这样可以防止出现 当配置文件中的name真的为nil时, 从内存中Get到的确实不是nil的情况。
		//         > 比如使用Get时, 获得的可能是name= map[0:<nil> 1:<nil>]
		audioPackageConfig.Viper.Set("audio_files."+arg.Sha256+".name", nil)

		// 查看name在内存中的值, 是否可配置文件一致(已检测一致)
		// fmt.Println("audio_files."+arg.Sha256+".name=", audioPackageConfig.GetValue("audio_files."+arg.Sha256+".name"))

		// 每次删除后, 都需要判断是否需要删除音频文件(此处的判断, 依赖前一行对name的nil设置, 否则可能会获得内存中与实际文件中不一致的值, 参考上方tips)
		if audioPackageConfig.GetValue("audio_files."+arg.Sha256+".name") == nil {

			audioPkgUUID, ok := audioPackageConfig.GetValue("audio_pkg_uuid").(string)
			if !ok {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 获取音频包UUID失败",
				})
				return
			}

			// 删除音频源文件
			err := os.Remove(filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", arg.Sha256+arg.Type))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 删除音频文件失败:" + err.Error(),
				})
				logger.Error("message", "error: 删除音频文件失败:"+err.Error())
				return
			}
			// 音频源文件删除成功后，删除配置项中的音频文件配置项
			audioPackageConfig.DeleteValue("audio_files." + arg.Sha256)
		}

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	keytonePkgRouters.POST("/play_sound", func(ctx *gin.Context) {

		type Arg struct {
			Sha256    string  `json:"sha256"`
			Type      string  `json:"type"`
			StartTime float64 `json:"startTime"`
			EndTime   float64 `json:"endTime"`
			Volume    float64 `json:"volume"`
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.Sha256 == "" || arg.Type == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		audioPkgUUID, ok := audioPackageConfig.GetValue("audio_pkg_uuid").(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取音频包UUID失败",
			})
			return
		}

		go keySound.PlayKeySound(&keySound.AudioFilePath{Part: filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", arg.Sha256+arg.Type)}, &keySound.Cut{
			StartMS: int64(arg.StartTime),
			EndMS:   int64(arg.EndTime),
			Volume:  arg.Volume,
		})

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	keytonePkgRouters.GET("/get_audio_package_list", func(ctx *gin.Context) {
		list, err := audioPackageList.GetAudioPackageList(audioPackageConfig.AudioPackagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取音频包列表失败:" + err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "ok",
			"list":    list,
		})
	})

	keytonePkgRouters.GET("/get_audio_package_name", func(ctx *gin.Context) {

		path := ctx.DefaultQuery("path", "unknown")
		if path == "unknown" || path == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容path值, 不符合接口规定格式:",
			})
			return
		}
		name := audioPackageList.GetAudioPackageName(path)
		ctx.JSON(200, gin.H{
			"message": "ok",
			"name":    name,
		})
	})

	keytonePkgRouters.POST("/export_album", func(ctx *gin.Context) {
		type Arg struct {
			AlbumPath string `json:"albumPath"`
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.AlbumPath == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值不符合接口规定格式:" + err.Error(),
			})
			return
		}

		// 检查源文件夹是否存在且可访问
		srcInfo, err := os.Stat(arg.AlbumPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 源专辑文件夹不存在或无法访问:" + err.Error(),
			})
			return
		}
		if !srcInfo.IsDir() {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 源路径不是一个文件夹",
			})
			return
		}

		// 使用内存buffer来创建zip
		buffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buffer)

		// 遍历键音专辑文件夹并添加到zip
		err = filepath.Walk(arg.AlbumPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("遍历文件夹失败: %v", err)
			}

			// 获取相对路径
			relPath, err := filepath.Rel(arg.AlbumPath, path)
			if err != nil {
				return fmt.Errorf("计算相对路径失败: %v", err)
			}

			// 如果是文件夹根目录,跳过
			if relPath == "." {
				return nil
			}

			// 创建zip文件头信息
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("创建文件头信息失败: %v", err)
			}
			header.Name = relPath

			if info.IsDir() {
				header.Name += "/" // 确保目录以/结尾
			} else {
				header.Method = zip.Deflate // 使用压缩
			}

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("创建zip条目失败: %v", err)
			}

			if info.IsDir() {
				return nil
			}

			// 以只读方式打开源文件
			file, err := os.OpenFile(path, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("打开源文件失败: %v", err)
			}
			defer file.Close()

			// 复制文件内容到zip
			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("写入zip文件失败: %v", err)
			}

			return nil
		})

		// 检查是否有压缩错误
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 压缩文件失败:" + err.Error(),
			})
			return
		}

		// 关闭zip writer
		if err = zipWriter.Close(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 关闭zip writer失败:" + err.Error(),
			})
			return
		}

		// 设置响应头并发送数据
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(arg.AlbumPath)))
		ctx.Data(http.StatusOK, "application/zip", buffer.Bytes())
	})

	keytonePkgRouters.POST("/delete_album", func(ctx *gin.Context) {
		type Arg struct {
			AlbumPath string `json:"albumPath"`
		}

		var arg Arg
		err := ctx.ShouldBind(&arg)
		if err != nil || arg.AlbumPath == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		err = os.RemoveAll(arg.AlbumPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 删除键音专辑失败:" + err.Error(),
			})
			return
		}

		// 清除sdk中的已选择键音包
		// TIPS: 若后续需要实现删除任意键音包, 则需要进行判断, 若删除的键音包中存在当前已选择的键音包, 才需要清除。
		audioPackageConfig.Viper = nil

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

}
