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

package keySound

import (
	"KeyTone/config"
	"KeyTone/logger"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	audioPackageConfig "KeyTone/audioPackage/config"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

//go:embed sounds
var sounds embed.FS

// 播放器的采样率为44100 Hz
var formatGlobalSampleRate beep.SampleRate = beep.SampleRate(44100)

// 在文件顶部添加变量
var (
	// 在文件顶部添加一个包级别的变量来存储循环播放的索引
	loopCurrentIndex map[string]uint
	// 添加一个静态变量来跟踪上一次的音频包UUID, 以便在切换音频包时重置loopCurrentIndex这个循环索引变量
	lastAudioPkgUUID string
)

func init() {

	// 初始化speaker。
	// 第二个参数的值, 不会对音质产生影响, 它只是缓冲区的大小。
	// > bufferSize参数指定扬声器缓冲区的样本数。更大的缓冲区大小意味着更低的CPU使用率和更可靠的播放。较低的缓冲区大小意味着更好的响应性和更少的延迟。
	// > * 缓冲区越大, cpu压力越小, 播放的整个过程崩溃率也会降低。(个人理解)
	// > * 缓冲区越小, cpu压力越大, 会得到更快的响应性和更少的延时。(个人理解)
	// > 鉴于个人的以上理解, 这个数值对我们KeyTone项目来说, 缓冲区设置的越小越好。
	// > * 但实际测试下来, 缓冲区无论如何设置, 其响应到播放完毕的用时都只有最大20ms作用的波动, 而且绝大部分时候, 波动仅有1ms左右。因此给其一个固定的值即可
	// starTime := time.Now()
	err := speaker.Init(formatGlobalSampleRate, formatGlobalSampleRate.N(time.Second/36))
	if err != nil {
		panic(err)
	}

	// 使用新的随机数生成器
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 仅开发环境下启用对 activeStreams 的监控打印功能
	if os.Getenv("SDK_MODE") == "debug" {
		go func() {
			for true {
				time.Sleep(3999 * time.Millisecond)
				PrintSyncMapDetailed(&activeStreams)
			}
		}()
	}
}

// 一个通用的可打印sync.Map内所有值的函数, 并在末尾统计总条目的数量(即打印共有多少个元素), Debug模式下使用。
func PrintSyncMapDetailed(m *sync.Map) {
	fmt.Println("===== sync.Map Detailed Contents =====")
	count := 0

	m.Range(func(key, value interface{}) bool {
		keyType := fmt.Sprintf("%T", key)
		valueType := fmt.Sprintf("%T", value)

		fmt.Printf("Key <%s>: %v\n", keyType, key)
		fmt.Printf("Value <%s>: %v\n", valueType, value)
		fmt.Println("---------------------------------")

		count++
		return true
	})

	fmt.Printf("Total entries: %d\n", count)
	fmt.Println("=====================================")
}

//region 设置管理所有正在播放的流的一系列变量及方法

// activeStreams 使用 sync.Map 来并发安全地存储当前所有正在播放的音频流。
// key 是指向 managedStream 的指针 (*managedStream)，它本身就是唯一的内存地址。
// value 是一个空结构体 struct{}{}, 不占用任何额外内存。
var activeStreams sync.Map
var activeStreamsAllDeleteFlag bool = false

// managedStream 包装了 beep.StreamSeekCloser 以便管理其生命周期。
// 它不需要 id 字段。
type managedStream struct {
	beep.StreamSeekCloser
}

// Stream 重写了内嵌的 Streamer 的 Stream 方法。
// 当流自然播放完毕时，它会用自身的指针作为键，将自己从 activeStreams 中删除。
func (ms *managedStream) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = ms.StreamSeekCloser.Stream(samples)
	if !ok {
		// 使用 ms (自身的指针)作为 key 来删除。
		activeStreams.Delete(ms)
	}
	return n, ok
}

// Close 重写了内嵌的 Closer 的 Close 方法。
// 它确保在关闭原始流之前，先将自己从 activeStreams 中删除。
func (ms *managedStream) Close() error {
	// 同样，使用 ms (自身的指针)作为 key 来删除。
	activeStreams.Delete(ms)
	return ms.StreamSeekCloser.Close()
}

// JoinManage 接收一个原始流，将其包装后进行管理。
// 它直接返回被包装后的流对象。
// 调用者使用返回的这个对象来进行播放和后续的控制。
//
// @param streamer: 原始的音频流。
// @return beep.StreamSeekCloser: 已被包装、可供播放和控制的流。
func JoinManage(streamer beep.StreamSeekCloser) beep.StreamSeekCloser {
	managed := &managedStream{
		StreamSeekCloser: streamer,
	}

	// 使用 managed 这个指针作为 key 存储，value 为空结构体。
	activeStreams.Store(managed, struct{}{})

	return managed
}

// CloseStream 根据传入的流对象引用来关闭它。
//
// @param stream: JoinManage 函数返回的那个流对象。
// @return bool: 如果成功找到并关闭了流，则返回 true；否则返回 false。
func CloseStream(stream beep.StreamCloser) bool {
	// key 就是 stream 对象本身。
	// 我们使用 Load 来检查它是否存在于 map 中。
	if _, ok := activeStreams.Load(stream); ok {
		// 存在，则调用它的 Close 方法。
		// 这将触发我们重写的 Close()，它会负责从 map 中删除。
		stream.Close()
		return true
	}
	// 如果流已经自然结束并从 map 中移除，Load 会失败，这里返回 false。
	return false
}

// CloseAllStreams 会关闭当前所有正在管理的音频流。
func CloseAllStreams() {
	activeStreamsAllDeleteFlag = true
	// time.Sleep(10 * time.Millisecond) // 等待所有将会触发的音频流进入activeStreams的管理范畴, 以保证100%全部关闭。
	time.Sleep(30 * time.Millisecond) // 虽然经上万次测试验证, 10ms的延时已足够, 但为防万一, 将此时间提高三倍。
	activeStreams.Range(func(key, value interface{}) bool {
		// 在这个设计中，key 就是我们要操作的流对象。
		if stream, ok := key.(beep.StreamCloser); ok {
			// 调用 Close 会触发包装器的 Close 方法，
			// 该方法会从 map 中安全地删除这个条目。
			stream.Close()
		}
		// 返回 true 以继续遍历并关闭所有流。
		return true
	})
	activeStreamsAllDeleteFlag = false
}

//endregion 设置管理所有正在播放的流的一系列变量及方法

type AudioFilePath struct {
	SS     string // 优先级最低
	Global string // 优先级仅次于Part
	Part   string // 优先级最高
}

type Cut struct {
	StartMS int64
	EndMS   int64 // 当 EndMS 小于或等于 StartMS  时, 不会播放任何声音
	Volume  float64
}

// 键音播放器
//
// Parameters:
//   - audioFilePath - 指定音频文件路径的结构体, 为nil代表不播放任何音频。
//   - cut - 裁剪键音的必要结构体, 为nil代表不裁剪。
//   - isPreviewMode - 可选参数, 用于指示是否为预览模式（使用原始音量）
//
// Returns:
//   - void
func PlayKeySound(audioFilePath *AudioFilePath, cut *Cut, isPreviewMode ...bool) {
	// 保证在删除全部活动流期间, 不新增任何播放项
	if activeStreamsAllDeleteFlag {
		return
	}

	if audioFilePath == nil {
		return
	}

	var audioFile fs.File
	var err error  // 注意, 这里一定要同时带上err。 否则在if else 内部, 和已声明的audioFile一起取返回值而临时创建的err, 会造成已声明的audioFile被重新声明并定义, 从而发生作用域问题。
	var ext string // 用于判断音频类型
	if audioFilePath.Part != "" {
		audioFile, err = os.Open(audioFilePath.Part)
		ext = strings.ToLower(filepath.Ext(audioFilePath.Part))
	} else if audioFilePath.Global != "" {
		audioFile, err = os.Open(audioFilePath.Global)
		ext = strings.ToLower(filepath.Ext(audioFilePath.Global))
	} else {
		audioFile, err = sounds.Open("sounds/" + audioFilePath.SS)
		ext = strings.ToLower(filepath.Ext(audioFilePath.SS))
	}

	if err != nil {
		logger.Error("message", fmt.Sprintf("error: failed to open audio file: %v", err))
		return
	}
	defer audioFile.Close()

	// 对文件进行解码
	audioStreamer, format, err := decodeAudioFile(audioFile, ext)
	if err != nil {
		logger.Error("message", fmt.Sprintf("error: failed to decode audio file: %v", err))
		return
	}
	audioStreamer = JoinManage(audioStreamer)
	defer audioStreamer.Close()

	fmt.Println("format.SampleRate", format.SampleRate)
	fmt.Println("formatGlobalSampleRate", formatGlobalSampleRate)

	// 将文件的采样率, 设置成与播放器一致
	reStreamer := beep.Resample(4, format.SampleRate, formatGlobalSampleRate, audioStreamer)

	// 处理cut参数
	endSample := -1 // 为保证cut=nil时, 也能正常保留原始工作。(当从配置文件获取的信息达不到构造cut时, cut将不会被构造。cut释放为nil的逻辑不应该在播放器端处理<如start和end都等于0时, cut就应该为nil, 即全量PlayKeySound播放>。)
	initVolume := 0.0
	// 如果cut=nil则全量播放
	if cut != nil {
		startSample := 0
		startSample = format.SampleRate.N(time.Millisecond * time.Duration(cut.StartMS))
		audioStreamer.Seek(startSample)
		// 若有不合理错误, 则直接退出, 不播放任何声音。
		// * 如果开始时间等于结束时间, 说明用户不想播放任何声音, 为避免内存浪费, 我们在此处也直接做退出处理。
		if cut.EndMS <= cut.StartMS {
			return
		}
		endSample = format.SampleRate.N(time.Millisecond * time.Duration(cut.EndMS))
		initVolume = cut.Volume
	}

	// 处理音量
	volume := &effects.Volume{
		Streamer: reStreamer,
		Base:     1.6,
		Volume:   initVolume,
		Silent:   false,
	}

	// 检查是否为预览模式（使用原始音量）
	shouldUseRawVolume := false
	if len(isPreviewMode) > 0 && isPreviewMode[0] {
		shouldUseRawVolume = true
	}

	// 仅在非预览模式时应用全局音量处理
	if !shouldUseRawVolume {
		volume = globalAudioVolumeAmplifyProcessing(volume)
		volume = globalAudioVolumeNormalProcessing(volume)
	}

	// ctrl := &beep.Ctrl{Streamer: volume, Paused: false}

	// 播放音乐
	done := make(chan struct{})
	defer close(done)
	// speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- struct{}{}
	})))

	// // FIXME: 暂时先如此处理, 后续再进一步处理
	// go (func() {
	// 	time.Sleep(500 * time.Millisecond)
	// 	ctrl.Paused = true
	// 	done <- true
	// })()

	// 等待播放完成
	re := true
	for re {
		select {
		case <-done:
			re = false
		case <-time.After(10 * time.Millisecond):
			pos := audioStreamer.Position()
			if pos >= endSample && endSample != -1 {
				// speaker.Lock()
				// ctrl.Paused = true // 目前只能用此一种方式, 在指定时间中止正在播放的音频 (由于暂停后, 会永远的滞留在播放器中等待恢复, 无法进入结束状态而被正确回收, 因此我们暂时采用静音的方式解决问题)
				// volume.Silent = true // 静音的方式解决问题, 虽然可以保证最终的内存正常释放, 但如果音频文件过大, 仍是会在一定时间内造成不必要的短暂内存泄漏问题。
				// volume.Silent = true // 仍保留这个的原因是: 为了防止末尾仍有声音, 或者说保证声音的纯净。
				// audioStreamer.Seek(audioStreamer.Len()) // 直接将其播放进度设置到末尾, 以使其直接播放完毕而自动调用内存回收。(从而避免音频文件过大时, 在一定时间内造成的短暂不必要的内存占用过大问题。)
				audioStreamer.Close()
				// speaker.Unlock()
				<-done // 为了防止beep.Callback回调卡死而造成的内存泄漏, 这里必须如此处理(就算提前结束, 也要正确的等待Callback回调)
				re = false
			}
		}
	}
	// fmt.Println("播放用时", time.Since(starTime))
	fmt.Println("结束------结束------结束")
}

func decodeAudioFile(file fs.File, ext string) (beep.StreamSeekCloser, beep.Format, error) {
	switch ext {
	case ".wav":
		return wav.Decode(file)
	case ".mp3":
		return mp3.Decode(file)
	case ".ogg":
		return vorbis.Decode(file)
	default:
		return nil, beep.Format{}, errors.New("unsupported audio format: " + ext)
	}
}

func globalAudioVolumeAmplifyProcessing(audioStreamer beep.Streamer) *effects.Volume {
	audio_volume_processing_volume_amplify, ok := config.GetValue("audio_volume_processing.volume_amplify").(float64)
	if !ok {
		audio_volume_processing_volume_amplify = config.Audio_volume_processing___volume_amplify
		go config.SetValue("audio_volume_processing.volume_amplify", config.Audio_volume_processing___volume_amplify)
	}

	return &effects.Volume{
		Streamer: audioStreamer,
		Base:     1.6,
		Volume:   audio_volume_processing_volume_amplify,
		Silent:   false,
	}

}

func globalAudioVolumeNormalProcessing(audioStreamer beep.Streamer) *effects.Volume {

	// 取出正常全局音量
	main_home_audio_volume_processing_volume_normal, ok := config.GetValue("main_home.audio_volume_processing.volume_normal").(float64)
	if !ok {
		main_home_audio_volume_processing_volume_normal = config.Main_home___audio_volume_processing___volume_normal
		go config.SetValue("main_home.audio_volume_processing.volume_normal", config.Main_home___audio_volume_processing___volume_normal)
	}

	// 最大音量为正常全局音量(0), 确保其不能超出正常的原始音量值
	// * FIXME:  这个可以在之后测试viper是否修复了内存覆盖bug
	//          * 1. 在安全范围内手动修改配置文件对应字段<此处为normal<=0>。
	//          * 2. 确认viper.WatchConfig()监听到的真实配置文件的对应字段的修改<此时可以由音量大小是否改变来判断 或 通过查看viper.Get()是否可以获取此字段的真实变化来判断>。
	//          * 3. 通过viper.Set()设置某个字段并使用viper.WriteConfig()写入配置文件后。<手动修改配置文件的normal>0即可自动触发此部分>
	//          * 4. 重复`步骤1` 和 `步骤2`,判断是否发生内存覆盖bug<此时可以由音量大小是否改变来判断 或 通过查看viper.Get()是否可以获取此字段的真实变化来判断>
	if main_home_audio_volume_processing_volume_normal > 0 {
		main_home_audio_volume_processing_volume_normal = 0.0
		go config.SetValue("main_home.audio_volume_processing.volume_normal", 0.0)
	}

	// 最小音量为 增强的音量-用户定义的正常音量缩减范围限度(-volume_amplify - volume_normal_reduce_scope), 确保其最小音量不能超过定义范围
	// * TODO: 等上方内存覆盖bug修复后,为了重新支持 通过修改配置文件 更改运行中程序设置的功能, 这里需要对最小限度进行判断
	// if main_home_audio_volume_processing_volume_normal < ???

	// 处理是否静音
	main_home_audio_volume_processing_volume_silent, ok := config.GetValue("main_home.audio_volume_processing.volume_silent").(bool)
	if !ok {
		main_home_audio_volume_processing_volume_silent = config.Main_home___audio_volume_processing___volume_silent
		go config.SetValue("main_home.audio_volume_processing.volume_silent", config.Main_home___audio_volume_processing___volume_silent)
	}

	return &effects.Volume{
		Streamer: audioStreamer,
		Base:     1.6,
		Volume:   main_home_audio_volume_processing_volume_normal,
		Silent:   main_home_audio_volume_processing_volume_silent,
	}

}

func init() {
	// initKeyDownSoundBuffer()
	// initKeyUpSoundBuffer()
}

var bufferKeyDownSound *beep.Buffer
var bufferKeyUpSound *beep.Buffer

func initKeyDownSoundBuffer() {
	audioFile, err := sounds.Open("sounds/down.wav")
	if err != nil {
		panic(err)
	}
	defer audioFile.Close()

	// 对文件进行解码
	audioStreamer, format, err := wav.Decode(audioFile)
	if err != nil {
		panic(err)
	}
	defer audioStreamer.Close()

	bufferKeyDownSound = beep.NewBuffer(format)
	bufferKeyDownSound.Append(audioStreamer)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/36))
}

func initKeyUpSoundBuffer() {
	audioFile, err := sounds.Open("sounds/up.wav")
	if err != nil {
		panic(err)
	}
	defer audioFile.Close()

	// 对文件进行解码
	audioStreamer, format, err := wav.Decode(audioFile)
	if err != nil {
		panic(err)
	}
	defer audioStreamer.Close()

	bufferKeyUpSound = beep.NewBuffer(format)
	bufferKeyUpSound.Append(audioStreamer)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/36))
}

func KeyDownSoundPlay() {
	shot := bufferKeyDownSound.Streamer(0, bufferKeyDownSound.Len())
	speaker.Play(shot)

	// // 播放音乐
	// done := make(chan bool)
	// speaker.Play(beep.Seq(shot, beep.Callback(func() {
	// 	done <- true
	// })))

	// // 等待播放完成
	// <-done
}

func KeyUpSoundPlay() {
	shot := bufferKeyUpSound.Streamer(0, bufferKeyUpSound.Len())
	speaker.Play(shot)

	// // 播放音乐
	// done := make(chan bool)
	// speaker.Play(beep.Seq(shot, beep.Callback(func() {
	// 	done <- true
	// })))

	// // 等待播放完成
	// <-done
}

// 添加常量定义按键状态
const (
	KeyStateDown = "down"
	KeyStateUp   = "up"
)

// ConfigGetter 用于抽象读取配置的函数（支持 editor 模式的 Viper 与 route 模式的只读快照）。
type ConfigGetter func(string) any

// getValue 对 ConfigGetter 进行统一封装，避免 nil 调用。
func getValue(get ConfigGetter, key string) any {
	if get == nil {
		return nil
	}
	return get(key)
}

// getAudioPkgUUID 从配置中读取 audio_pkg_uuid，并在必要时回退到目录名。
func getAudioPkgUUID(get ConfigGetter, fallback string) (string, bool) {
	if get == nil {
		return "", false
	}
	value, ok := get("audio_pkg_uuid").(string)
	if ok && strings.TrimSpace(value) != "" {
		return value, true
	}
	if fallback != "" {
		return fallback, true
	}
	return "", false
}

func resolvePlaybackConfig(keycode string) (ConfigGetter, string) {
	// =============================
	// 选择播放来源（核心逻辑）
	// =============================
	//
	// 约定：KeyTone 将鼠标按键编码为负数（例如 "-1" ~ "-5"），键盘按键为非负字符串。
	// 因此可用 keycode 是否以 "-" 开头来区分鼠标事件。
	//
	// 模式：
	// - editor：编辑试听模式。播放读取 audioPackageConfig.Viper（可写配置，可能伴随 SSE 更新）。
	// - route-unified / route-split：主页路由模式。播放读取只读快照（AlbumSnapshot），由 apply_playback_routing 预加载。
	//
	// 热路径约束：本函数 MUST 只做内存读取；不得触发磁盘 IO / 解密 / JSON 解析。
	// 这些昂贵操作必须提前在 ApplyPlaybackRouting 中完成。
	state := GetPlaybackState()
	isMouse := strings.HasPrefix(keycode, "-")

	if state.SourceMode == SourceModeEditor {
		// 编辑试听模式下直接使用当前可编辑配置。
		// 注意：audioPackageConfig.Viper 由 LoadConfig 初始化；若为空表示尚未加载专辑。
		if audioPackageConfig.Viper == nil {
			return nil, ""
		}
		uuid, _ := getAudioPkgUUID(audioPackageConfig.GetValue, "")
		return audioPackageConfig.GetValue, uuid
	}

	if state.SourceMode == SourceModeRouteUnified || state.SourceMode == SourceModeRouteSplit {
		var snapshot *AlbumSnapshot
		if state.SourceMode == SourceModeRouteUnified {
			// 统一路由：键盘/鼠标共享快照。
			snapshot = state.Routing.UnifiedSnapshot
		} else {
			// =============================
			// 分离路由：键盘/鼠标各自使用独立快照
			// =============================
			if isMouse {
				snapshot = state.Routing.MouseSnapshot
				// 回退逻辑：仅当用户在设置中启用 mouse_fallback_to_keyboard 时，
				// 鼠标专辑缺失才会回退到键盘专辑。
				// 默认行为：彻底分离，鼠标无专辑则无声（返回 nil 后使用内嵌测试音或静音）。
				if snapshot == nil {
					// 读取回退开关配置（热路径内只做内存读取，配置已在 SDK 初始化时加载）
					fallbackEnabled, _ := config.GetValue("playback.routing.mouse_fallback_to_keyboard").(bool)
					if fallbackEnabled {
						snapshot = state.Routing.KeyboardSnapshot
					}
				}
			} else {
				snapshot = state.Routing.KeyboardSnapshot
			}
		}
		if snapshot != nil && snapshot.Viper != nil {
			return snapshot.GetValue, snapshot.AudioPkgUUID()
		}
	}

	// 未命中任何来源时回退到内嵌测试音。
	return nil, ""
}

// 音频包处理器
// * 此函数会根据处理结果来调用播放器播放对应的音频结果。
func KeySoundHandler(keyState string, keycode string) {
	configGetter, audioPkgUUID := resolvePlaybackConfig(keycode)

	// audioPkgUUID 的用途：当配置引用 audio_files（sha256+ext）时，需要拼出真实文件路径：
	//   AudioPackagePath/<audioPkgUUID>/audioFiles/<sha256><ext>
	// 在路由快照里，我们优先读取配置内的 audio_pkg_uuid；若缺失则回退为专辑目录名。

	// 如果没有选择音频包，则默认使用内嵌的测试音频进行播放
	if configGetter == nil {
		switch keycode {
		case "-1", "-2", "-3", "-4", "-5":
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_down.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: 0.6})
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_up.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: -0.6})
			}
		default:
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "test_down.MP3",
				}, &Cut{StartMS: 32, EndMS: 100, Volume: 0})
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "test_up.MP3",
				}, &Cut{StartMS: 28, EndMS: 100, Volume: 0})
			}
		}

		return
	}
	// 从音频包配置中获取相关设置, 并根据配置决定如何播放

	// TODO: 根据传入的具体按键Keycode, 来独立寻找其预设的播放配置, 以播放对应音频。
	single := getValue(configGetter, "key_tone.single."+keycode+"."+keyState)
	fmt.Println("single====", single)
	if single != nil {
		// 将single转换为map类型以便访问其中的值
		soundEffectType := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".type")
		// TIPS: 这个虽然 single和global都有, 但也没必要提取, 因为它仍旧只会执行一次。(但提取后, 会使得仅播放嵌入测试音时, 也执行这个无关紧要的逻辑)
		audioPkgUUID, ok := audioPkgUUID, audioPkgUUID != ""
		if !ok {
			logger.Error("message", "error: 获取音频包UUID失败")
			return
		}

		// TIPS: 没必要将single和global的 handleSoundEffect 的逻辑抽离到一个函数内。 因为这样我们在改传参的基础上, 还需要改返回值 并在此处调用后 通过返回值判断是否return。
		if soundEffectType == "audio_files" {
			sha256, ok := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".value.sha256").(string)
			if !ok {
				return
			}
			fileType, ok := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".value.type").(string)
			if !ok {
				return
			}
			audio_file_name := sha256 + fileType
			audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)

			PlayKeySound(&AudioFilePath{
				Global: audio_file_path,
			}, nil)
			return
		}

		if soundEffectType == "sounds" {
			sound_UUID := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".value")
			if sound_UUID == nil {
				return
			}

			soundParsePlayWith(configGetter, sound_UUID.(string), audioPkgUUID)

			return
		}

		if soundEffectType == "key_sounds" {
			key_sound_UUID := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".value")
			if key_sound_UUID == nil {
				return
			}

			keySoundParsePlayWith(configGetter, key_sound_UUID.(string), keyState, audioPkgUUID, true, keycode, 0)
			// PlayKeySound(&AudioFilePath{
			// 	Global: audio_file_path,
			// }, nil)
			return
		}

	}

	// TODO: 若具体按键配置为空, 则根据全局配置决定如何播放
	global := getValue(configGetter, "key_tone.global."+keyState)
	fmt.Println("global====", global)
	// * 如果global不为空, 则根据global的值来决定如何播放, 否则使用后续逻辑中的默认音频
	if global != nil {
		// 将global转换为map类型以便访问其中的值
		soundEffectType := getValue(configGetter, "key_tone.global."+keyState+".type")
		// soundEffectValue := audioPackageConfig.GetValue("key_tone.global." + keyState + ".value")
		audioPkgUUID, ok := audioPkgUUID, audioPkgUUID != ""
		if !ok {
			logger.Error("message", "error: 获取音频包UUID失败")
			return
		}

		if soundEffectType == "audio_files" {
			sha256, ok := getValue(configGetter, "key_tone.global."+keyState+".value.sha256").(string)
			if !ok {
				return
			}
			fileType, ok := getValue(configGetter, "key_tone.global."+keyState+".value.type").(string)
			if !ok {
				return
			}
			audio_file_name := sha256 + fileType
			audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)

			PlayKeySound(&AudioFilePath{
				Global: audio_file_path,
			}, nil)
			return
		}

		if soundEffectType == "sounds" {
			sound_UUID := getValue(configGetter, "key_tone.global."+keyState+".value")
			if sound_UUID == nil {
				return
			}

			soundParsePlayWith(configGetter, sound_UUID.(string), audioPkgUUID)

			return
		}

		if soundEffectType == "key_sounds" {
			key_sound_UUID := getValue(configGetter, "key_tone.global."+keyState+".value")
			if key_sound_UUID == nil {
				return
			}

			keySoundParsePlayWith(configGetter, key_sound_UUID.(string), keyState, audioPkgUUID, true, keycode, 0)
			// PlayKeySound(&AudioFilePath{
			// 	Global: audio_file_path,
			// }, nil)
			return
		}

	}

	// 若全局配置中为空, 则获取配置中内置测试音效的启用状态, 以决定是否使用默认音频进行播放。(优先级最低)
	// * 我们没有对is_enable_embedded_test_sound做类型断言, 因此其可能为nil或bool,
	isEnableEmbeddedTestSound := getValue(configGetter, "key_tone.is_enable_embedded_test_sound."+keyState)
	// * 只要不是主动设置为false, 我们都使用默认音频
	if isEnableEmbeddedTestSound == true || isEnableEmbeddedTestSound == nil {
		switch keycode {
		case "-1", "-2", "-3", "-4", "-5":
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_down.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: 0.6})
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_up.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: -0.6})
			}
		default:
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "test_down.MP3",
				}, &Cut{StartMS: 32, EndMS: 100, Volume: 0})
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "test_up.MP3",
				}, &Cut{StartMS: 28, EndMS: 100, Volume: 0})
			}
		}
		return
	}

}

// 声音解析, 获取 实际音频文件的路径 以及 裁剪的参数
// 参数:
//   - sound_UUID: 声音的UUID值
//   - audioPkgUUID: 音频包的UUID
func soundParsePlay(sound_UUID string, audioPkgUUID string) {
	soundParsePlayWith(audioPackageConfig.GetValue, sound_UUID, audioPkgUUID)
}

func soundParsePlayWith(get ConfigGetter, sound_UUID string, audioPkgUUID string) {
	sha256, ok := getValue(get, "sounds."+sound_UUID+".source_file_for_sound"+".sha256").(string)
	if !ok {
		logger.Error("message", "error: sha256 value is nil or not a string")
		return
	}

	fileType, ok := getValue(get, "sounds."+sound_UUID+".source_file_for_sound"+".type").(string)
	if !ok {
		logger.Error("message", "error: file type value is nil or not a string")
		return
	}

	audio_file_name := sha256 + fileType
	audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
	cut := &Cut{
		StartMS: int64(getValue(get, "sounds."+sound_UUID+".cut.start_time").(float64)),
		EndMS:   int64(getValue(get, "sounds."+sound_UUID+".cut.end_time").(float64)),
		Volume:  getValue(get, "sounds."+sound_UUID+".cut.volume").(float64),
	}
	PlayKeySound(&AudioFilePath{
		Global: audio_file_path,
	}, cut)
}

// 键音解析, 获取 实际音频文件的路径 以及 播放参数
// 参数:
//   - key_sound_UUID: 至臻键音的UUID索引值
//   - keyState: 按键状态 (如 "down", "up" 等)
//   - audioPkgUUID: 音频包的UUID
//   - isGlobal: 是否为全局音效
//   - keycode: 按键码
//
// 功能说明:
//
// 1. 根据mode判断播放模式:
//   - single: 单一音效模式,按顺序播放配置的音效
//   - random: 随机音效模式,随机选择一个音效播放
//   - loop:   循环音效模式,循环播放配置的音效
func keySoundParsePlay(key_sound_UUID string, keyState string, audioPkgUUID string, isGlobal bool, keycode string, count uint16) {
	keySoundParsePlayWith(audioPackageConfig.GetValue, key_sound_UUID, keyState, audioPkgUUID, isGlobal, keycode, count)
}

func keySoundParsePlayWith(get ConfigGetter, key_sound_UUID string, keyState string, audioPkgUUID string, isGlobal bool, keycode string, count uint16) {
	// 此处限制键音的嵌套数量上限为1000层, 这样即使键音专辑中存在至臻键音间的无限循环依赖也不必担心因此可能引起的  cpu超负荷风险 或 内存占用过多的内存溢出风险。
	// * 理论上, 设置9999甚至更高也是可行的, 但没必要, 因为没有人会去制作继承嵌套超过1000层的键音。
	if count > 1000 {
		return
	}
	count = count + 1

	mode := getValue(get, "key_sounds."+key_sound_UUID+"."+keyState+".mode")
	if mode == "single" {
		value := getValue(get, "key_sounds."+key_sound_UUID+"."+keyState+".value")

		if value != nil {
			for _, v := range value.([]interface{}) {
				vMap := v.(map[string]interface{})
				if vMap["type"] == "audio_files" {
					valueMap := vMap["value"].(map[string]interface{})
					audio_file_name := valueMap["sha256"].(string) + valueMap["type"].(string)
					audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
					PlayKeySound(&AudioFilePath{
						Global: audio_file_path,
					}, nil)
					return
				}
						if vMap["type"] == "sounds" {
							sound_UUID := vMap["value"].(string)
							soundParsePlayWith(get, sound_UUID, audioPkgUUID)
					return
				}
						if vMap["type"] == "key_sounds" {
							key_sound_UUID := vMap["value"].(string)
							keySoundParsePlayWith(get, key_sound_UUID, keyState, audioPkgUUID, isGlobal, keycode, count)
					return
				}
			}
		}
	}
	if mode == "random" {
		value := getValue(get, "key_sounds."+key_sound_UUID+"."+keyState+".value")
		if value != nil {
			values := value.([]interface{})
			// TIPS: 防止因空值造成后续步骤panic。
			if len(values) == 0 {
				return
			}
			// 创建一个新的随机数生成器实例
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			randomIndex := r.Intn(len(values))
			logger.Debug("随机算法检测", "randomIndex====", randomIndex)
			v := values[randomIndex]
			vMap := v.(map[string]interface{})

			if vMap["type"] == "audio_files" {
				valueMap := vMap["value"].(map[string]interface{})
				audio_file_name := valueMap["sha256"].(string) + valueMap["type"].(string)
				audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
				PlayKeySound(&AudioFilePath{
					Global: audio_file_path,
				}, nil)
				return
			}
			if vMap["type"] == "sounds" {
				sound_UUID := vMap["value"].(string)
				soundParsePlayWith(get, sound_UUID, audioPkgUUID)
				return
			}
			if vMap["type"] == "key_sounds" {
				key_sound_UUID := vMap["value"].(string)
				keySoundParsePlayWith(get, key_sound_UUID, keyState, audioPkgUUID, isGlobal, keycode, count)
				return
			}
		}
	}
	if mode == "loop" {
		value := getValue(get, "key_sounds."+key_sound_UUID+"."+keyState+".value")
		if value != nil {

			// 检测音频包是否发生切换
			if lastAudioPkgUUID != audioPkgUUID {
				loopCurrentIndex = make(map[string]uint) // 重置循环计数器
				lastAudioPkgUUID = audioPkgUUID
			}

			values := value.([]interface{})
			// TIPS: 防止因空值造成后续步骤panic。
			if len(values) == 0 {
				return
			}

			// 使用简单的key来标识不同的键音
			var key string
			if isGlobal {
				key = fmt.Sprintf("global_%s_%s_%s", keycode, key_sound_UUID, keyState)
			} else {
				key = fmt.Sprintf("%s_%s_%s", keycode, key_sound_UUID, keyState)
			}

			// 获取当前索引
			currentIndex := loopCurrentIndex[key]

			// 如果当前索引超出数组长度，重置为0
			if currentIndex >= uint(len(values)) {
				currentIndex = 0
			}

			// 获取当前要播放的值
			v := values[currentIndex]
			vMap := v.(map[string]interface{})

			// 更新索引，简单加1即可
			loopCurrentIndex[key] = currentIndex + 1

			// 根据类型播放音频
			if vMap["type"] == "audio_files" {
				valueMap := vMap["value"].(map[string]interface{})
				audio_file_name := valueMap["sha256"].(string) + valueMap["type"].(string)
				audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
				PlayKeySound(&AudioFilePath{
					Global: audio_file_path,
				}, nil)
				return
			}
			if vMap["type"] == "sounds" {
				sound_UUID := vMap["value"].(string)
				soundParsePlayWith(get, sound_UUID, audioPkgUUID)
				return
			}
			if vMap["type"] == "key_sounds" {
				key_sound_UUID := vMap["value"].(string)
				keySoundParsePlayWith(get, key_sound_UUID, keyState, audioPkgUUID, isGlobal, keycode, count)
				return
			}
		}
	}
}
