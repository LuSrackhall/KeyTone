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
	"KeyTone/audioPackage/enc"
	audioPackageList "KeyTone/audioPackage/list"
	"KeyTone/config"
	"KeyTone/keyEvent"
	"KeyTone/keySound"
	"KeyTone/logger"
	"KeyTone/signature"
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	KeytoneMagicNumber   = "KTAF"    // KeyTone Album Format
	KeytoneVersion       = "1.0.0"   // 当前版本号
	KeytoneFileSignature = "KTALBUM" // 文件签名
	KeytoneFileVersion   = 1         // 文件版本（已废弃，仅用于向后兼容）
)

// ==============================
// 专辑导出文件对称密钥（版本化）
// 注意：这些变量不再是 const，而是 var，以便在编译时通过 -ldflags 进行注入
// 注入的值应为经过 XOR 混淆后的 Hex 字符串（与授权流一致）
// ==============================

// xorMask 用于混淆密钥的掩码，必须与授权流一致
var xorMask = []byte{0x55, 0xAA, 0x33, 0xCC, 0x99, 0x66, 0x11, 0xEE, 0x77, 0xBB, 0x22, 0xDD, 0x88, 0x44, 0xFF, 0x00}

// 默认开源密钥常量（明文）
const (
	DefaultKeytoneEncryptKeyV1      = "KeyTone2024SecretKey"                  // v1 密钥（旧版本，用于向后兼容）
	DefaultKeytoneEncryptKeyV2      = "KeyTone2025AlbumSecureEncryptionKeyV2" // v2 密钥（当前版本）
	DefaultKeytoneEncryptKeyCurrent = DefaultKeytoneEncryptKeyV2
)

// 版本化加密密钥（可注入）
var (
	KeytoneEncryptKeyV1      = DefaultKeytoneEncryptKeyV1
	KeytoneEncryptKeyV2      = DefaultKeytoneEncryptKeyV2
	KeytoneEncryptKeyCurrent = DefaultKeytoneEncryptKeyCurrent
	KeytoneEncryptKey        = KeytoneEncryptKeyV1 // 已废弃：向后兼容，请使用 KeytoneEncryptKeyV1
)

func deobfuscateString(obfuscatedHex string) string {
	obfuscated, err := hex.DecodeString(obfuscatedHex)
	if err != nil {
		// 非 hex（可能是默认明文，或用户错误注入了明文）
		return obfuscatedHex
	}
	realBytes := make([]byte, len(obfuscated))
	for i, b := range obfuscated {
		realBytes[i] = b ^ xorMask[i%len(xorMask)]
	}
	return string(realBytes)
}

func getPlainEncryptKey(value string, defaultValue string) string {
	if value == defaultValue {
		return defaultValue
	}
	return deobfuscateString(value)
}

// KeytoneAlbumMeta 用于存储专辑元数据
type KeytoneAlbumMeta struct {
	MagicNumber string    `json:"magicNumber"`
	Version     string    `json:"version"`
	ExportTime  time.Time `json:"exportTime"`
	AlbumUUID   string    `json:"albumUUID"`
	AlbumName   string    `json:"albumName"`
}

type KeyStateMessage struct {
	Type    string `json:"type"`
	Keycode uint16 `json:"keycode"`
	State   string `json:"state"`
}

// KeytoneFileHeader 文件头结构
type KeytoneFileHeader struct {
	Signature [7]byte  // "KTALBUM"
	Version   uint8    // 文件版本
	DataSize  uint64   // 加密后的zip数据大小
	Checksum  [32]byte // zip数据的SHA-256校验和
}

// 验证 nanoid 格式的辅助函数
func isValidNanoID(id string) bool {
	// nanoid 默认使用 21 个字符，字符集为 A-Za-z0-9_-
	if len(id) != 21 {
		return false
	}
	for _, char := range id {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return true
}

// 验证专辑结构的辅助函数
func isValidAlbumStructure(albumPath string) error {
	// 检查目录名是否符合 nanoid 格式
	dirName := filepath.Base(albumPath)
	if !isValidNanoID(dirName) {
		return fmt.Errorf("专辑目录名不符合规范")
	}

	// 检查 package.json || config.json 是否存在
	configPath := filepath.Join(albumPath, "package.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(albumPath, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return fmt.Errorf("缺少必要的配置文件 package.json 或 config.json")
		}
	}

	return nil
}

// 将 copyDir 函数移到 package server 级别
func copyDir(src string, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// 复制文件
			srcFile, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}

			// 保持文件权限
			srcInfo, err := os.Stat(srcPath)
			if err != nil {
				return err
			}
			err = os.Chmod(dstPath, srcInfo.Mode())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 将 validateAlbumMeta 函数移到包级别，和其他辅助函数放在一起
func validateAlbumMeta(zipReader *zip.ReadCloser) error {
	// 查找元数据文件
	var metaFile *zip.File
	for _, f := range zipReader.File {
		if f.Name == ".keytone-album" {
			metaFile = f
			break
		}
	}

	if metaFile == nil {
		return fmt.Errorf("不是有效的 KeyTone 专辑文件：缺少元数据")
	}

	// 读取元数据文件
	rc, err := metaFile.Open()
	if err != nil {
		return fmt.Errorf("读取元数据失败: %v", err)
	}
	defer rc.Close()

	var meta KeytoneAlbumMeta
	if err := json.NewDecoder(rc).Decode(&meta); err != nil {
		return fmt.Errorf("解析元数据失败: %v", err)
	}

	// 验证魔数
	if meta.MagicNumber != KeytoneMagicNumber {
		return fmt.Errorf("不是有效的 KeyTone 专辑文件：无效的魔数")
	}

	// 验证版本兼容性
	// 这里可以添加版本兼容性检查的逻辑

	return nil
}

// 处理导入文件的通用函数
func processImportedFile(src io.Reader, header KeytoneFileHeader, tempZipPath string) error {
	// 读取加密的数据
	encryptedData := make([]byte, header.DataSize)
	if _, err := io.ReadFull(src, encryptedData); err != nil {
		return &ImportError{Message: "读取文件数据失败:" + err.Error()}
	}

	// 解密数据（支持多版本密钥）
	zipData, usedVersion, err := decryptAlbumData(encryptedData, header)
	if err != nil {
		return &ImportError{Message: err.Error()}
	}
	logger.Info("专辑数据解密成功", "file_version", header.Version, "used_key_version", usedVersion)

	// 保存解密后的数据到临时zip文件
	if err := os.WriteFile(tempZipPath, zipData, 0644); err != nil {
		return &ImportError{Message: "保存临时文件失败:" + err.Error()}
	}

	return nil
}

// 解压并验证专辑结构的通用函数
func extractAndValidateAlbum(zipReader *zip.ReadCloser, tempDir string) (string, error) {
	// 解压到临时目录
	for _, file := range zipReader.File {
		// 构建完整的目标路径
		targetPath := filepath.Join(tempDir, file.Name)

		// 确保目标目录存在
		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		// 创建目标文件的父目录
		os.MkdirAll(filepath.Dir(targetPath), 0755)

		// 创建目标文件
		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", fmt.Errorf("创建目标文件失败: %v", err)
		}

		// 打开源文件
		inFile, err := file.Open()
		if err != nil {
			outFile.Close()
			return "", fmt.Errorf("打开源文件失败: %v", err)
		}

		// 复制文件内容
		_, err = io.Copy(outFile, inFile)
		outFile.Close()
		inFile.Close()
		if err != nil {
			return "", fmt.Errorf("复制文件内容失败: %v", err)
		}
	}

	// 获取解压后的专辑目录
	files, err := os.ReadDir(tempDir)
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("读取解压目录失败或目录为空: %v", err)
	}

	var albumDir os.DirEntry
	for _, f := range files {
		if f.IsDir() {
			if albumDir != nil {
				return "", fmt.Errorf("zip 文件中包含多个目录")
			}
			albumDir = f
		}
	}

	if albumDir == nil {
		return "", fmt.Errorf("zip 文件中未找到专辑目录")
	}

	albumPath := filepath.Join(tempDir, albumDir.Name())
	if err := isValidAlbumStructure(albumPath); err != nil {
		return "", fmt.Errorf("无效的专辑格式: %v", err)
	}

	return albumPath, nil
}

// 导入错误类型
type ImportError struct {
	Message string
}

func (e *ImportError) Error() string {
	return e.Message
}

// 简单的异或加密/解密函数（对称加密）
func xorCrypt(data []byte, key string) []byte {
	keyBytes := []byte(key)
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}
	return result
}

// getEncryptKeyByVersion 根据版本号获取对应的加密密钥
// 参数:
//   - version: 文件版本号
//
// 返回:
//   - string: 对应的加密密钥
func getEncryptKeyByVersion(version uint8) string {
	switch version {
	case 1:
		return getPlainEncryptKey(KeytoneEncryptKeyV1, DefaultKeytoneEncryptKeyV1)
	case 2:
		return getPlainEncryptKey(KeytoneEncryptKeyV2, DefaultKeytoneEncryptKeyV2)
	default:
		// 未知版本，返回当前密钥
		logger.Warn("未知的文件版本号，使用当前密钥", "version", version)
		return getPlainEncryptKey(KeytoneEncryptKeyCurrent, DefaultKeytoneEncryptKeyCurrent)
	}
}

// decryptAlbumData 解密专辑数据，支持多版本密钥回退
// 参数:
//   - encryptedData: 加密的数据
//   - header: 文件头结构
//
// 返回:
//   - []byte: 解密后的数据
//   - uint8: 实际使用的密钥版本
//   - error: 错误信息
func decryptAlbumData(encryptedData []byte, header KeytoneFileHeader) ([]byte, uint8, error) {
	// 首先根据版本号选择密钥
	decryptKey := getEncryptKeyByVersion(header.Version)
	zipData := xorCrypt(encryptedData, decryptKey)

	// 验证校验和
	checksum := sha256.Sum256(zipData)
	if checksum == header.Checksum {
		return zipData, header.Version, nil
	}

	// 如果校验失败且版本不是v1，尝试使用v1密钥回退
	if header.Version != 1 {
		logger.Warn("使用版本密钥解密失败，尝试v1密钥回退", "version", header.Version)
		zipData = xorCrypt(encryptedData, getEncryptKeyByVersion(1))
		checksum = sha256.Sum256(zipData)
		if checksum == header.Checksum {
			logger.Info("使用v1密钥成功解密", "file_version", header.Version)
			return zipData, 1, nil
		}
	}

	return nil, 0, fmt.Errorf("文件校验失败，文件可能已损坏或使用了不支持的加密版本")
}

func ServerRun() {

	// 启动签名名片图片清理任务（在SDK启动5秒后执行一次）
	go func() {
		time.Sleep(5 * time.Second)
		// CleanupOrphanCardImages 的参数为历史兼容保留；实际解密逻辑使用 KeyA（支持构建注入）
		if err := signature.CleanupOrphanCardImages(signature.GetKeyA()); err != nil {
			logger.Error("签名名片图片清理任务执行失败", "error", err.Error())
		}
	}()

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
	signatureRouters(r)

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

	// 创建一个channel用于服务器就绪通知
	ready := make(chan bool, 1)

	// 使用listener启动服务
	go func() {
		// 启动服务器
		go func() {
			// time.Sleep(10000 * time.Millisecond)
			if err := r.RunListener(listener); err != nil {
				logger.Error("服务器启动失败:", err)
				ready <- false
			}
		}()

		for {
			// 这里我们没有设置超时限制, 所以会一直阻塞等待所请求的相关服务的返回信息, 直到返回成功或失败, 并解除阻塞。(未设置超时限制, 不会返回超时信息)(由于不做超时限制, 会节省一些损耗, 也能第够一时间作出响应)
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", port))
			// 如果请求没有出错测继续。(否则会开启一轮新的请求->这里我不确定err的可能情况, 因为err不为nil的几率小到可以忽略不计->不知道之前设置的超时限制触发的是不是这里的err, 几率小到离谱, 懒得测试了, 到此为止)。
			if err == nil {
				resp.Body.Close()
				// 如果请求成功则向通道发送true, 以向终端输出端口号信息。(如果失败则不做任何处理, 让本grouting自行结束即可)(如果失败的话, 相关服务启动失败后就会向通道发送false了->以向终端输出服务启动失败的相关信息, 故此处无需处理。)
				if resp.StatusCode == 200 {
					ready <- true
					// fmt.Println("55555555555555666666666666666666666666666666777")
					return
				}
				// fmt.Println("55555555555555")
				// 只要请求本身没有出错, 就退出循环, 不再进行重新请求。
				return
			}
			// fmt.Println("55555555555555666666666666666666666666666666")
		}
	}()

	// 等待服务器就绪信号
	isReady := <-ready
	if !isReady {
		fmt.Println("SDK的本地server模块启动失败")
		return
	}
	// 输出端口信息，让Electron主进程可以捕获
	fmt.Printf("KEYTONE_PORT=%d\n", port)
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

	// 应用播放路由（只读快照加载）
	keytonePkgRouters.POST("/apply_playback_routing", func(ctx *gin.Context) {
		type Arg struct {
			Mode             string `json:"mode"`
			UnifiedAlbumPath string `json:"unifiedAlbumPath"`
			KeyboardAlbumPath string `json:"keyboardAlbumPath"`
			MouseAlbumPath   string `json:"mouseAlbumPath"`
		}

		// 请求语义：
		// - mode=unified：只使用 unifiedAlbumPath（键盘/鼠标共享）
		// - mode=split：使用 keyboardAlbumPath + mouseAlbumPath（键盘/鼠标分离）
		// - album path 既可传“绝对路径”，也可仅传“UUID”（SDK 会拼成 AudioPackagePath/UUID）

		var arg Arg
		if err := ctx.ShouldBindJSON(&arg); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		// 该调用会在内存中生成只读快照（viper.New() + ReadConfig），保证播放热路径无磁盘 IO。
		// 返回值 result 内包含每个来源的 requested/resolved/loaded/error，便于前端诊断。
		result, err := keySound.ApplyPlaybackRouting(arg.Mode, arg.UnifiedAlbumPath, arg.KeyboardAlbumPath, arg.MouseAlbumPath)
		if err != nil {
			// partial：允许“部分成功”（例如键盘快照成功但鼠标失败）。
			// 前端应根据 result 判断是否需要提示用户、或回退到统一模式。
			ctx.JSON(200, gin.H{
				"message": "partial",
				"error":   err.Error(),
				"result":  result,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "ok",
			"result":  result,
		})
	})

	// 切换播放来源模式（editor / route）
	keytonePkgRouters.POST("/set_playback_source_mode", func(ctx *gin.Context) {
		type Arg struct {
			Mode            string `json:"mode"`
			EditorAlbumPath string `json:"editorAlbumPath"`
		}

		var arg Arg
		if err := ctx.ShouldBindJSON(&arg); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收--收到的前端数据内容值, 不符合接口规定格式:" + err.Error(),
			})
			return
		}

		// editor：播放依赖可编辑配置（audioPackageConfig.Viper），用于编辑页实时试听。
		// route：播放依赖只读快照（由 apply_playback_routing 生成），用于主页日常播放。
		state, err := keySound.SetPlaybackSourceMode(arg.Mode, arg.EditorAlbumPath)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "ok",
			"mode":    state.SourceMode,
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

		audioPackageConfig.SetValue("audio_files."+arg.Sha256+".name", nil)
		temp := audioPackageConfig.GetValue("audio_files." + arg.Sha256 + ".name")
		if m, ok := temp.(map[string]any); ok {
			if len(m) <= 1 {
				// 删除音频文件, 因为此时音频文件的存在毫无意义
				audioPkgUUID, ok := audioPackageConfig.GetValue("audio_pkg_uuid").(string)
				if !ok {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "error: 获取音频包UUID失败",
					})
					return
				}

				// 在正式删除音频源文件之前, 需要先释放所有流的文件句柄, 因为在Win系统中, 不释放的话是没办法成功关闭的。
				keySound.CloseAllStreams()
				// time.Sleep(10 * time.Millisecond)
				// // 需要调用两次的原因是 -> 前端在单击ui中的删除按钮时的行为本身, 会增加一个额外的正在播放的声音流, 而由于sync.map天然的锁机制, 它并不会包含在上述的关闭流程中。
				// keySound.CloseAllStreams() // 由于CloseAllStreams()函数内部已通过升级变得足够可靠, 因此无需再进行二次调用。

				// 删除音频源文件
				err := os.Remove(filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", arg.Sha256+arg.Type))
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "error: 删除音频文件失败:" + err.Error(),
					})
					logger.Error("message", "error: 删除音频文件失败:"+err.Error())
					return
				}

				// 音频源文件删除成功后，删除配置项中的音频文件配置项(此时不需要管具体的NameID的删除, 因为我们已经从父级删除了)
				audioPackageConfig.DeleteValue("audio_files." + arg.Sha256)

				ctx.JSON(200, gin.H{
					"message": "ok",
				})

			} else {
				// 仅删除对应的NameID, 因为还有其它的NameID需要依赖相关的音频文件。
				audioPackageConfig.DeleteValue("audio_files." + arg.Sha256 + ".name." + arg.NameID)
				ctx.JSON(200, gin.H{
					"message": "ok",
				})

			}
		} else {
			// 类型断言失败，说明 temp 不是 map[string]interface{} 类型
			// 可以处理错误或忽略
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 删除音频文件失败——断言失败",
			})
		}

	})

	keytonePkgRouters.POST("/play_sound", func(ctx *gin.Context) {

		type Arg struct {
			Sha256        string  `json:"sha256"`
			Type          string  `json:"type"`
			StartTime     float64 `json:"startTime"`
			EndTime       float64 `json:"endTime"`
			Volume        float64 `json:"volume"`
			IsPreviewMode bool    `json:"isPreviewMode"`
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
		}, "", arg.IsPreviewMode)

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// ============================================================================
	// GET /get_audio_stream - 读取当前编辑专辑下的音频源文件（用于前端波形渲染）
	//
	// 说明：
	// - 前端通过 sha256 + type（扩展名，如 ".wav"）定位 audioFiles/<sha256><type>
	// - 该接口返回音频文件内容，支持 Range（由 http.ServeContent 处理）
	// - 仅用于编辑器场景（依赖 audio_pkg_uuid 已加载）
	// ============================================================================
	keytonePkgRouters.GET("/get_audio_stream", func(ctx *gin.Context) {
		type Arg struct {
			Sha256 string `form:"sha256"`
			Type   string `form:"type"`
		}

		var arg Arg
		if err := ctx.ShouldBindQuery(&arg); err != nil || arg.Sha256 == "" || arg.Type == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"message": "error: 参数接收失败",
			})
			return
		}

		// 基础校验：防止路径穿越/非法扩展名
		if len(arg.Sha256) != 64 {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "error: sha256 格式不正确"})
			return
		}
		if _, err := hex.DecodeString(arg.Sha256); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "error: sha256 格式不正确"})
			return
		}
		if !strings.HasPrefix(arg.Type, ".") {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "error: type 格式不正确"})
			return
		}
		for _, r := range arg.Type[1:] {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
				ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "error: type 格式不正确"})
				return
			}
		}

		audioPkgUUID, ok := audioPackageConfig.GetValue("audio_pkg_uuid").(string)
		if !ok || audioPkgUUID == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取音频包UUID失败",
			})
			return
		}

		fileName := arg.Sha256 + arg.Type
		filePath := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", fileName)

		f, err := os.Open(filePath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "error: 音频文件不存在"})
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error: 读取音频文件失败"})
			return
		}

		// content-type：根据扩展名推断
		contentType := mime.TypeByExtension(arg.Type)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		ctx.Header("Content-Type", contentType)
		ctx.Header("Content-Disposition", "inline; filename=\""+fileName+"\"")

		http.ServeContent(ctx.Writer, ctx.Request, fileName, info.ModTime(), f)
	})

	// ============================================================================
	// GET /get_audio_package_list - 获取专辑列表（含签名摘要）
	//
	// 功能说明：
	//   - 获取所有专辑的路径列表
	//   - 同时获取每个专辑的签名摘要信息（直接导出作者名称、图片）
	//   - 签名摘要用于前端列表展示，避免逐个请求签名详情
	//
	// 返回格式：
	//   {
	//     "message": "ok",
	//     "list": ["path1", "path2", ...],
	//     "signatureInfo": {
	//       "path1": { "hasSignature": true, "directExportAuthorName": "...", "directExportAuthorImage": "..." },
	//       "path2": { "hasSignature": false }
	//     }
	//   }
	// ============================================================================
	keytonePkgRouters.GET("/get_audio_package_list", func(ctx *gin.Context) {
		list, err := audioPackageList.GetAudioPackageList(audioPackageConfig.AudioPackagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取音频包列表失败:" + err.Error(),
			})
			return
		}

		// 获取每个专辑的签名摘要信息
		// TIPS: 遍历过程中为每个专辑独立创建和释放 Viper 实例，避免内存泄漏
		signatureInfo := make(map[string]*audioPackageList.AlbumSignatureSummary)
		for _, albumPath := range list {
			summary, err := audioPackageList.GetAlbumSignatureSummary(albumPath)
			if err != nil {
				// 获取签名摘要失败时，使用默认值（无签名）
				logger.Warn("获取专辑签名摘要失败", "albumPath", albumPath, "err", err.Error())
				signatureInfo[albumPath] = &audioPackageList.AlbumSignatureSummary{HasSignature: false}
			} else {
				signatureInfo[albumPath] = summary
			}
		}

		ctx.JSON(200, gin.H{
			"message":       "ok",
			"list":          list,
			"signatureInfo": signatureInfo,
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
		albumName := audioPackageList.GetAudioPackageName(path)
		retryCount := 0
		for albumName == nil {
			albumName = audioPackageList.GetAudioPackageName(path)
			// 添加一个计数器和最大重试次数(防止死循环造成的资源占用问题)
			if retryCount >= 6 {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 获取专辑名称失败: 超过最大重试次数",
				})
				return
			}
			retryCount++
			fmt.Println("GetAudioPackageName为nil, 尝试重新获取。 获取次数=", retryCount)
			time.Sleep(100 * time.Millisecond) // 添加短暂延迟
		}
		albumName, ok := albumName.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取专辑名称失败: 类型转换错误",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "ok",
			"name":    albumName,
		})
	})

	// 独立加密 API：前端在选择"需要签名"时调用，根据传入的专辑路径加密配置文件
	keytonePkgRouters.POST("/encrypt_album_config", func(ctx *gin.Context) {
		type Arg struct {
			AlbumPath string `json:"albumPath"`
		}

		var arg Arg
		if err := ctx.ShouldBindJSON(&arg); err != nil || strings.TrimSpace(arg.AlbumPath) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 请求体必须包含 albumPath",
			})
			return
		}

		albumPath := arg.AlbumPath
		// 校验目录是否存在
		info, err := os.Stat(albumPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 专辑目录不存在或无法访问:" + err.Error(),
			})
			return
		}
		if !info.IsDir() {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: albumPath 必须是目录",
			})
			return
		}

		pkgPath := filepath.Join(albumPath, "package.json")
		stubInfo, pkgRaw, err := audioPackageConfig.ReadCoreStubInfo(albumPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 读取指示 JSON 失败:" + err.Error(),
			})
			return
		}
		if stubInfo != nil {
			corePath := filepath.Join(albumPath, stubInfo.Core)
			if _, err := os.Stat(corePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: core 文件缺失:" + err.Error(),
				})
				return
			}
			logger.Info("专辑配置已加密，跳过重复操作", "album", albumPath)
			ctx.JSON(http.StatusOK, gin.H{
				"message":           "ok",
				"already_encrypted": true,
			})
			return
		}

		if pkgRaw == nil {
			pkgRaw, err = os.ReadFile(pkgPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 读取配置文件失败:" + err.Error(),
				})
				return
			}
		}

		albumUUID := filepath.Base(albumPath)
		plainJSON := strings.TrimSpace(string(pkgRaw))
		migrated := false
		if enc.IsLikelyHexCipher(pkgRaw) {
			decoded, decErr := enc.DecryptConfigHex(plainJSON, albumUUID)
			if decErr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 旧版密文解密失败:" + decErr.Error(),
				})
				return
			}
			plainJSON = decoded
			migrated = true
		} else {
			if err := enc.ValidateJSONFast(plainJSON); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "error: 配置内容不是有效的 JSON",
				})
				return
			}
		}

		cipherBytes, err := enc.EncryptConfigBytes(plainJSON, albumUUID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 加密配置失败:" + err.Error(),
			})
			return
		}

		corePath := filepath.Join(albumPath, audioPackageConfig.CoreFileName)
		tmpCore := corePath + ".tmp"
		if err := os.WriteFile(tmpCore, cipherBytes, 0644); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 写入 core 临时文件失败:" + err.Error(),
			})
			return
		}
		if err := os.Rename(tmpCore, corePath); err != nil {
			_ = os.Remove(tmpCore)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 重命名 core 临时文件失败:" + err.Error(),
			})
			return
		}

		if err := audioPackageConfig.WriteCoreStubFile(albumPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 写入指示 JSON 失败:" + err.Error(),
			})
			return
		}

		logger.Info("专辑配置加密成功", "album", albumPath, "migrated", migrated)
		ctx.JSON(http.StatusOK, gin.H{
			"message":   "ok",
			"encrypted": true,
			"migrated":  migrated,
		})

		audioPackageConfig.LoadConfig(albumPath, false)
		ctx.JSON(http.StatusOK, gin.H{
			"message":   "ok",
			"encrypted": true,
		})
	})

	keytonePkgRouters.POST("/apply_signature_config", func(ctx *gin.Context) {
		var req struct {
			AlbumPath              string `json:"albumPath" binding:"required"`
			NeedSignature          bool   `json:"needSignature"`
			RequireAuthorization   bool   `json:"requireAuthorization"`
			SignatureID            string `json:"signatureId" binding:"required"`
			ContactEmail           string `json:"contactEmail"`
			ContactAdditional      string `json:"contactAdditional"`
			UpdateSignatureContent bool   `json:"updateSignatureContent"`
			// AuthorizationUUID 授权标识UUID
			// 首次导出时由前端nanoid生成并传入，用于未来签名授权导出/导入功能
			// 再次导出时可传空字符串，SDK会沿用已存储的UUID
			AuthorizationUUID string `json:"authorizationUUID"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的签名配置请求体",
			})
			return
		}

		if strings.TrimSpace(req.AlbumPath) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: albumPath 不能为空",
			})
			return
		}

		if (req.NeedSignature || req.RequireAuthorization) && strings.TrimSpace(req.SignatureID) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 选择签名后才能导出",
			})
			return
		}

		if req.RequireAuthorization && strings.TrimSpace(req.ContactEmail) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 需要授权时必须提供联系邮箱",
			})
			return
		}

		logger.Info("收到签名配置应用请求",
			"albumPath", req.AlbumPath,
			"needSignature", req.NeedSignature,
			"requireAuthorization", req.RequireAuthorization,
			"signatureId", req.SignatureID,
			"contactEmail", req.ContactEmail,
			"contactAdditional", req.ContactAdditional,
			"authorizationUUID", req.AuthorizationUUID,
		)

		// 调用签名应用函数
		qualificationCode, err := audioPackageConfig.ApplySignatureToAlbum(
			req.AlbumPath,
			req.SignatureID,
			req.RequireAuthorization,
			req.ContactEmail,
			req.ContactAdditional,
			req.UpdateSignatureContent,
			req.AuthorizationUUID, // 授权标识UUID（首次导出时前端nanoid生成）
		)

		if err != nil {
			logger.Error("签名应用失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 签名应用失败: " + err.Error(),
			})
			return
		}

		logger.Info("签名应用成功", "qualificationCode", qualificationCode)

		ctx.JSON(http.StatusOK, gin.H{
			"message":           "ok",
			"success":           true,
			"qualificationCode": qualificationCode,
		})
	})

	// 获取专辑签名信息（前端需求2和4）
	keytonePkgRouters.POST("/get_album_signature_info", func(ctx *gin.Context) {
		var req struct {
			AlbumPath string `json:"albumPath" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的请求参数",
			})
			return
		}

		signatureInfo, err := audioPackageConfig.GetAlbumSignatureInfo(req.AlbumPath)
		if err != nil {
			logger.Error("获取专辑签名信息失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    signatureInfo,
		})
	})

	// 读取专辑内文件（用于展示签名图片等）
	keytonePkgRouters.POST("/get_album_file", func(ctx *gin.Context) {
		var req struct {
			AlbumPath    string `json:"albumPath" binding:"required"`
			RelativePath string `json:"relativePath" binding:"required"` // 相对于专辑目录的路径，如 "audioFiles/xxx.jpg"
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的请求参数",
			})
			return
		}

		// 构建完整文件路径
		fullPath := filepath.Join(req.AlbumPath, req.RelativePath)

		// 安全检查：确保请求的文件在专辑目录内
		cleanAlbumPath := filepath.Clean(req.AlbumPath)
		cleanFullPath := filepath.Clean(fullPath)
		if !strings.HasPrefix(cleanFullPath, cleanAlbumPath) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 非法的文件路径",
			})
			return
		}

		// 读取文件
		fileData, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("读取专辑内文件失败", "path", fullPath, "error", err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "error: 文件不存在或无法读取",
			})
			return
		}

		// 检测文件类型
		contentType := http.DetectContentType(fileData)

		// 返回文件内容
		ctx.Data(http.StatusOK, contentType, fileData)
	})

	// 检查签名是否在专辑中（前端需求3）
	keytonePkgRouters.POST("/check_signature_in_album", func(ctx *gin.Context) {
		var req struct {
			AlbumPath   string `json:"albumPath" binding:"required"`
			SignatureID string `json:"signatureId" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的请求参数",
			})
			return
		}

		isInAlbum, qualCode, hasChanges, err := audioPackageConfig.CheckSignatureInAlbum(req.AlbumPath, req.SignatureID)
		if err != nil {
			logger.Error("检查签名失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":           "ok",
			"isInAlbum":         isInAlbum,
			"qualificationCode": qualCode,
			"hasChanges":        hasChanges,
		})
	})

	// 检查签名授权状态（前端需求3）
	keytonePkgRouters.POST("/check_signature_authorization", func(ctx *gin.Context) {
		var req struct {
			AlbumPath   string `json:"albumPath" binding:"required"`
			SignatureID string `json:"signatureId" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的请求参数",
			})
			return
		}

		isAuthorized, requireAuth, qualCode, err := audioPackageConfig.CheckSignatureAuthorization(req.AlbumPath, req.SignatureID)
		if err != nil {
			logger.Error("检查签名授权失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":              "ok",
			"isAuthorized":         isAuthorized,
			"requireAuthorization": requireAuth,
			"qualificationCode":    qualCode,
		})
	})

	// 获取可用于导出的签名列表（前端需求3）
	keytonePkgRouters.POST("/get_available_signatures", func(ctx *gin.Context) {
		var req struct {
			AlbumPath string `json:"albumPath" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的请求参数",
			})
			return
		}

		signatures, err := audioPackageConfig.GetAvailableSignaturesForExport(req.AlbumPath)
		if err != nil {
			logger.Error("获取可用签名列表失败", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":    "ok",
			"signatures": signatures,
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

		// 创建 zip buffer
		buffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buffer)

		// 创建并写入元数据文件
		metaWriter, err := zipWriter.Create(".keytone-album")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 创建元数据文件失败:" + err.Error(),
			})
			return
		}

		albumName, ok := audioPackageList.GetAudioPackageName(arg.AlbumPath).(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 获取专辑名称失败: 类型转换错误",
			})
			return
		}

		meta := KeytoneAlbumMeta{
			MagicNumber: KeytoneMagicNumber,
			Version:     KeytoneVersion,
			ExportTime:  time.Now(),
			AlbumUUID:   filepath.Base(arg.AlbumPath),
			AlbumName:   albumName,
		}

		metaJson, err := json.Marshal(meta)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 生成元数据失败:" + err.Error(),
			})
			return
		}

		_, err = metaWriter.Write(metaJson)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 写入元数据失败:" + err.Error(),
			})
			return
		}

		// 遍历键音专辑文件夹并添加到zip
		err = filepath.Walk(arg.AlbumPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("遍历文件夹失败: %v", err)
			}

			// 获取相对路径
			relPath, err := filepath.Rel(filepath.Dir(arg.AlbumPath), path)
			if err != nil {
				return fmt.Errorf("计算相对路径失败: %v", err)
			}

			// 统一使用正斜杠，确保跨平台兼容性
			relPath = filepath.ToSlash(relPath)

			// 创建zip文件头信息
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("创建文件头信息失败: %v", err)
			}

			// 使用标准化的路径
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

		// 获取 zip 数据
		zipData := buffer.Bytes()

		// 计算校验和
		checksum := sha256.Sum256(zipData)

		// 加密 zip 数据（使用当前版本密钥 v2）
		encryptedData := xorCrypt(zipData, getEncryptKeyByVersion(2))

		// 创建文件头（使用版本号 2）
		header := KeytoneFileHeader{
			Version:  2, // 使用 v2 版本号
			DataSize: uint64(len(encryptedData)),
			Checksum: checksum,
		}
		copy(header.Signature[:], KeytoneFileSignature)

		// 写入最终文件
		finalBuffer := new(bytes.Buffer)
		if err := binary.Write(finalBuffer, binary.LittleEndian, header); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 写入文件头失败:" + err.Error(),
			})
			return
		}
		finalBuffer.Write(encryptedData)

		// 设置响应头并发送数据
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.ktalbum", filepath.Base(arg.AlbumPath)))
		ctx.Data(http.StatusOK, "application/octet-stream", finalBuffer.Bytes())
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

		// * 在正式删除音频源文件之前, 需要先释放所有流的文件句柄, 因为在Win系统中, 不释放的话是没办法成功关闭的。
		keySound.CloseAllStreams()

		// * 正式删除现有目录
		err = os.RemoveAll(arg.AlbumPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 删除键音专辑失败:" + err.Error(),
			})
			return
		}

		// 清除sdk中的已选择键音包
		// TIPS: 若后续需要实现删除任意键音包, 则需要进行判断, 若删除的键音包中存在当前已选择的键音包, 才需要清除。
		audioPackageConfig.ClearConfig()

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// 修改导入处理函数
	// 导入为新专辑
	keytonePkgRouters.POST("/import_album_as_new", func(ctx *gin.Context) {
		// 获取上传的文件
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件上传失败:" + err.Error(),
			})
			return
		}

		// 获取新的专辑ID
		newAlbumId := ctx.PostForm("newAlbumId")
		if !isValidNanoID(newAlbumId) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的专辑ID格式",
			})
			return
		}

		// 检查文件扩展名
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".ktalbum") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式，请选择 .ktalbum 文件",
			})
			return
		}

		// 读取文件数据
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开文件失败:" + err.Error(),
			})
			return
		}
		defer src.Close()

		// 读取文件头并验证
		var header KeytoneFileHeader
		if err := binary.Read(src, binary.LittleEndian, &header); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 读取文件头失败:" + err.Error(),
			})
			return
		}

		// 验证文件签名
		if string(header.Signature[:]) != KeytoneFileSignature {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式：不是 KeyTone 专辑文件",
			})
			return
		}

		// 解压和处理文件
		tempDir, err := os.MkdirTemp("", "keytone_import_*")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 创建临时目录失败:" + err.Error(),
			})
			return
		}
		defer os.RemoveAll(tempDir)

		tempZipPath := filepath.Join(tempDir, "temp.zip")
		if err := processImportedFile(src, header, tempZipPath); err != nil {
			if importErr, ok := err.(*ImportError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "error: " + importErr.Message,
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		// 打开zip文件进行验证
		zipReader, err := zip.OpenReader(tempZipPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开zip文件失败:" + err.Error(),
			})
			return
		}
		defer zipReader.Close()

		// 验证zip内的元数据
		if err := validateAlbumMeta(zipReader); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		// 解压到临时目录并验证结构
		albumPath, err := extractAndValidateAlbum(zipReader, tempDir)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// 使用新的专辑ID创建目标路径
		originalAlbumID := filepath.Base(albumPath)
		targetPath := filepath.Join(audioPackageConfig.AudioPackagePath, newAlbumId)

		// 复制到目标路径
		if err := copyDir(albumPath, targetPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 复制专辑文件夹失败:" + err.Error(),
			})
			return
		}

		// 更新配置文件中的UUID
		if err := audioPackageList.UpdateAlbumUUID(targetPath, newAlbumId, originalAlbumID); err != nil {
			// 如果更新失败，清理已复制的文件夹
			os.RemoveAll(targetPath)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 更新专辑配置失败:" + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// 原有的导入专辑路由
	keytonePkgRouters.POST("/import_album", func(ctx *gin.Context) {
		// 获取上传的文件
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件上传失败:" + err.Error(),
			})
			return
		}

		// 检查是否是覆盖模式
		overwrite := ctx.PostForm("overwrite") == "true"

		// 检查文件扩展名
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".ktalbum") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式，请选择 .ktalbum 文件",
			})
			return
		}

		// 读取文件数据
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开文件失败:" + err.Error(),
			})
			return
		}
		defer src.Close()

		// 读取文件头
		var header KeytoneFileHeader
		if err := binary.Read(src, binary.LittleEndian, &header); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 读取文件头失败:" + err.Error(),
			})
			return
		}

		// 验证文件签名
		if string(header.Signature[:]) != KeytoneFileSignature {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式：不是 KeyTone 专辑文件",
			})
			return
		}

		// 读取加密的数据
		encryptedData := make([]byte, header.DataSize)
		if _, err := io.ReadFull(src, encryptedData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 读取文件数据失败:" + err.Error(),
			})
			return
		}

		// 解密数据（支持多版本密钥）
		zipData, usedVersion, err := decryptAlbumData(encryptedData, header)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}
		logger.Info("专辑数据解密成功", "file_version", header.Version, "used_key_version", usedVersion)

		// 创建临时文件保存 zip 数据
		tempDir, err := os.MkdirTemp("", "keytone_import_*")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 创建临时目录失败:" + err.Error(),
			})
			return
		}
		defer os.RemoveAll(tempDir)

		tempZipPath := filepath.Join(tempDir, "temp.zip")
		if err := os.WriteFile(tempZipPath, zipData, 0644); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 保存临时文件失败:" + err.Error(),
			})
			return
		}

		// 打开zip文件进行验证
		zipReader, err := zip.OpenReader(tempZipPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开zip文件失败:" + err.Error(),
			})
			return
		}
		defer zipReader.Close()

		// 验证zip内的元数据
		if err := validateAlbumMeta(zipReader); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}

		// 解压到临时目录
		for _, file := range zipReader.File {
			// 构建完整的目标路径
			targetPath := filepath.Join(tempDir, file.Name)

			// 确保目标目录存在
			if file.FileInfo().IsDir() {
				os.MkdirAll(targetPath, 0755)
				continue
			}

			// 创建目标文件的父目录
			os.MkdirAll(filepath.Dir(targetPath), 0755)

			// 创建目标文件
			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 创建目标文件失败:" + err.Error(),
				})
				return
			}

			// 打开源文件
			inFile, err := file.Open()
			if err != nil {
				outFile.Close()
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 打开源文件失败:" + err.Error(),
				})
				return
			}

			// 复制文件内容
			_, err = io.Copy(outFile, inFile)
			outFile.Close()
			inFile.Close()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 复制文件内容失败:" + err.Error(),
				})
				return
			}
		}

		// 获取解压后的专辑目录名
		files, err := os.ReadDir(tempDir)
		if err != nil || len(files) == 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 读取解压目录失败或目录为空:" + err.Error(),
			})
			return
		}

		// 检查是否只有一个目录
		var albumDir os.DirEntry
		for _, f := range files {
			if f.IsDir() {
				fmt.Println("跨平台一致性校验", f)
				if albumDir != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "error: zip 文件中包含多个目录",
					})
					return
				}
				albumDir = f
			}
		}

		if albumDir == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: zip 文件中未找到专辑目录",
			})
			return
		}

		albumPath := filepath.Join(tempDir, albumDir.Name())
		targetPath := filepath.Join(audioPackageConfig.AudioPackagePath, albumDir.Name())

		// 验证专辑结构
		if err := isValidAlbumStructure(albumPath); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的专辑格式: " + err.Error(),
			})
			return
		}

		// 获取新的专辑ID（如果提供）
		newAlbumId := ctx.PostForm("newAlbumId")
		if newAlbumId != "" {
			// 验证新的专辑ID是否符合nanoid格式
			if !isValidNanoID(newAlbumId) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "error: 无效的专辑ID格式",
				})
				return
			}
			// 使用新的专辑ID更新目标路径
			targetPath = filepath.Join(audioPackageConfig.AudioPackagePath, newAlbumId)
		}

		// 在复制到目标路径前检查是否存在
		if _, err := os.Stat(targetPath); err == nil {
			// 目标已存在
			if !overwrite {
				// 如果不是覆盖模式，返回特殊状态
				ctx.JSON(http.StatusOK, gin.H{
					"message": "album_exists",
				})
				return
			}
			// 覆盖模式：删除现有目录
			// * 在正式删除音频源文件之前, 需要先释放所有流的文件句柄, 因为在Win系统中, 不释放的话是没办法成功关闭的。
			keySound.CloseAllStreams()
			// time.Sleep(10 * time.Millisecond)
			// // * 需要调用两次的原因是 -> 前端在单击ui中的删除按钮时的行为本身, 会增加一个额外的正在播放的声音流, 而由于sync.map天然的锁机制, 它并不会包含在上述的关闭流程中。
			// keySound.CloseAllStreams() // 由于CloseAllStreams()函数内部已通过升级变得足够可靠, 因此无需再进行二次调用。

			// * 正式删除现有目录
			if err := os.RemoveAll(targetPath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "error: 删除现有专辑失败:" + err.Error(),
				})
				return
			}
		}

		// 使用复制替代移动
		if err := copyDir(albumPath, targetPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 复制专辑文件夹失败:" + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// 获取专辑文件的元数据信息
	keytonePkgRouters.POST("/get_album_meta", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 文件上传失败:" + err.Error(),
			})
			return
		}

		// 检查文件扩展名
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".ktalbum") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式，请选择 .ktalbum 文件",
			})
			return
		}

		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开文件失败:" + err.Error(),
			})
			return
		}
		defer src.Close()

		// 读取文件头
		var header KeytoneFileHeader
		if err := binary.Read(src, binary.LittleEndian, &header); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 读取文件头失败:" + err.Error(),
			})
			return
		}

		// 验证文件签名
		if string(header.Signature[:]) != KeytoneFileSignature {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 无效的文件格式：不是 KeyTone 专辑文件",
			})
			return
		}

		// 读取加密的数据
		encryptedData := make([]byte, header.DataSize)
		if _, err := io.ReadFull(src, encryptedData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 读取文件数据失败:" + err.Error(),
			})
			return
		}

		// 解密数据（支持多版本密钥）
		zipData, usedVersion, err := decryptAlbumData(encryptedData, header)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: " + err.Error(),
			})
			return
		}
		logger.Info("专辑数据解密成功", "file_version", header.Version, "used_key_version", usedVersion)

		// 创建临时zip文件
		tempFile, err := os.CreateTemp("", "keytone_meta_*.zip")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 创建临时文件失败:" + err.Error(),
			})
			return
		}
		tempPath := tempFile.Name()
		defer os.Remove(tempPath)
		defer tempFile.Close()

		// 写入zip数据
		if _, err := tempFile.Write(zipData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 写入临时文件失败:" + err.Error(),
			})
			return
		}
		tempFile.Close()

		// 打开zip文件
		zipReader, err := zip.OpenReader(tempPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 打开zip文件失败:" + err.Error(),
			})
			return
		}
		defer zipReader.Close()

		// 查找并读取元数据文件
		var metaFile *zip.File
		for _, f := range zipReader.File {
			if f.Name == ".keytone-album" {
				metaFile = f
				break
			}
		}

		if metaFile == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error: 不是有效的 KeyTone 专辑文件：缺少元数据",
			})
			return
		}

		// 读取元数据文件
		rc, err := metaFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 读取元数据失败:" + err.Error(),
			})
			return
		}
		defer rc.Close()

		var meta KeytoneAlbumMeta
		if err := json.NewDecoder(rc).Decode(&meta); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error: 解析元数据失败:" + err.Error(),
			})
			return
		}

		// 返回专辑元数据
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"meta":    meta,
		})
	})

}
