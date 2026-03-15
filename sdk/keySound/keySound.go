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

// errEmptyAudioCut 表示“这次裁剪在逻辑上不应该播放任何声音”。
// 这不是异常故障, 而是一个明确的业务分支:
// 1. 用户把结束时间拖到了开始时间之前/相同位置;
// 2. 裁剪区间完全落在音频总长度之外;
// 3. 音频长度本身异常, 无法得到任何可播放样本。
//
// 之所以单独定义这个错误, 是为了在调用方区分:
// - 真正需要记录日志的故障(如 Seek 失败、解码器异常);
// - 合法但不应发声的“空片段”情况。
var errEmptyAudioCut = errors.New("audio cut does not contain playable samples")

// preparePlaybackSource 将“裁剪描述 Cut”转换成一个真正可播放的 Streamer。
//
// 这是这次修复的核心: 旧实现是先 Resample, 再在播放期间轮询 Position(),
// 一旦发现到达结束点就直接 Close 原始流。这个做法的问题在于:
// 1. Resampler 为了插值会预读未来的一段原始样本;
// 2. 因此“原始流当前位置”并不等于“已经真正播出的声音位置”;
// 3. 尤其在很短的裁剪片段中, 预读进来的区间外样本会被插值带入输出;
// 4. 最终表现就是: 明明选中的片段没有声音, 实际却还能听到片段外的声音。
//
// 新方案改成两步:
// 1. 先在原始采样率的 StreamSeekCloser 上 Seek 到 startSample;
// 2. 再用 beep.Take 严格限制最多只能读取 end-start 个样本;
//
// 这样重采样器拿到的输入源本身就是一个“已经裁好长度”的只读片段,
// 它即便预读, 也只能在这个受限区间内预读, 不可能再越界读到片段外的声音。
//
// 返回值:
// - beep.Streamer: 供后续重采样与播放使用的源流;
// - float64: 初始音量偏移, 直接继承自 cut.Volume;
// - error: 区分空片段(errEmptyAudioCut)与真正故障。
func preparePlaybackSource(audioStreamer beep.StreamSeekCloser, sampleRate beep.SampleRate, cut *Cut) (beep.Streamer, float64, error) {
	initVolume := 0.0
	// 没有 cut 代表“播放整段音频”, 这时直接把原始流返回即可。
	if cut == nil {
		return audioStreamer, initVolume, nil
	}

	// 结束时间小于等于开始时间时, 语义上就是空区间, 必须明确无声返回。
	if cut.EndMS <= cut.StartMS {
		return nil, 0, errEmptyAudioCut
	}

	// 配置中的裁剪时间单位是毫秒, 这里统一转换为解码后“原始采样率”下的样本索引。
	// 注意必须使用原始采样率, 不能使用全局播放采样率, 否则时间轴会错位。
	startSample := sampleRate.N(time.Millisecond * time.Duration(cut.StartMS))
	endSample := sampleRate.N(time.Millisecond * time.Duration(cut.EndMS))
	totalSamples := audioStreamer.Len()

	// 没有任何可用样本时, 直接视为空片段。
	if totalSamples <= 0 {
		return nil, 0, errEmptyAudioCut
	}
	// 负值裁剪时间统一钳制到 0, 防止配置异常导致 Seek 到负位置。
	if startSample < 0 {
		startSample = 0
	}
	if endSample < 0 {
		endSample = 0
	}
	// 如果起点已经在文件末尾或更后面, 就没有任何可播放内容。
	if startSample >= totalSamples {
		return nil, 0, errEmptyAudioCut
	}
	// 结束位置允许越界, 但要钳制到文件真实长度, 等价于“播放到文件结尾”。
	if endSample > totalSamples {
		endSample = totalSamples
	}
	// 钳制后若仍然没有有效区间, 说明最终结果仍是空片段。
	if endSample <= startSample {
		return nil, 0, errEmptyAudioCut
	}

	// Seek 是必须检查错误的。
	// 旧逻辑忽略了 Seek 返回值, 一旦解码器拒绝 Seek 或位置异常, 播放可能回退到文件开头,
	// 这正是“选中的是静音段, 却听到别处声音”的另一个来源。
	if err := audioStreamer.Seek(startSample); err != nil {
		return nil, 0, fmt.Errorf("seek start sample %d failed: %w", startSample, err)
	}

	initVolume = cut.Volume
	// Take 会把源流严格截断为指定样本数。
	// 后续即便交给 Resample, Resample 也只能在这个已裁好的窗口中读取数据,
	// 不会再接触到区间外的原始样本。
	return beep.Take(endSample-startSample, audioStreamer), initVolume, nil
}

// 键音播放器
//
// Parameters:
//   - audioFilePath - 指定音频文件路径的结构体, 为nil代表不播放任何音频。
//   - cut - 裁剪键音的必要结构体, 为nil代表不裁剪。
//   - keycode - 当前按键/鼠标事件的 keycode（用于分离模式下的独立音量处理）
//   - isPreviewMode - 可选参数, 用于指示是否为预览模式（使用原始音量）
//
// Returns:
//   - void
func PlayKeySound(audioFilePath *AudioFilePath, cut *Cut, keycode string, keyState string, isPreviewMode ...bool) {
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

	// 先把“文件 + cut 配置”整理成一个真正可播放的源流。
	// 这里的关键原则是: 先裁剪, 再重采样。
	// 如果顺序反过来, 重采样器内部的预读行为就可能把裁剪区间外的内容提前读进来。
	playbackSource, initVolume, err := preparePlaybackSource(audioStreamer, format.SampleRate, cut)
	if err != nil {
		// 空片段不记错误日志, 直接静默返回, 这是符合用户配置语义的结果。
		if errors.Is(err, errEmptyAudioCut) {
			return
		}
		// 只有真正的异常才记录日志, 便于后续排查具体音频文件或解码器问题。
		logger.Error("message", fmt.Sprintf("error: failed to prepare playback source: %v", err))
		return
	}

	// 将文件的采样率, 设置成与播放器一致。
	// 裁剪必须先作用在原始采样率的流上, 再交给重采样器, 否则重采样器的预读会把裁剪区间外的数据带入播放。
	reStreamer := beep.Resample(4, format.SampleRate, formatGlobalSampleRate, playbackSource)

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
		// 分离模式下的键盘/鼠标独立音量叠加（在全局音量基础上再次叠加）
		volume = splitRouteAudioVolumeNormalProcessing(volume, keycode)
		// 按下/抬起音量单独控制叠加（默认关闭）
		volume = pressReleaseAudioVolumeNormalProcessing(volume, keycode, keyState)
		// 随机音量叠加（默认关闭，且仅做衰减）
		volume = randomAudioVolumeProcessing(volume)
		// 按下/抬起随机音量单独控制叠加（默认关闭，可与全局随机音量叠加）
		volume = pressReleaseRandomAudioVolumeProcessing(volume, keycode, keyState)
	}

	// ctrl := &beep.Ctrl{Streamer: volume, Paused: false}

	// 播放音乐
	// 这里使用一个带 1 个缓冲的 done 通道等待播放完成:
	// 1. speaker.Play 是异步的, 不会阻塞当前 goroutine;
	// 2. beep.Callback 会在整段流真正播放结束后触发;
	// 3. 带缓冲是为了防止极端调度下, 回调先于接收方执行时发生阻塞;
	// 4. select + default 则保证回调至多投递一次完成信号。
	//
	// 旧实现之所以复杂, 是因为它依赖“播放途中主动关流”来截断, 所以必须轮询、抢时机、担心回调卡死。
	// 现在裁剪已经由 Take 在数据源层面完成, 这里就只需要等待自然播放结束即可。
	done := make(chan struct{}, 1)
	// speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		select {
		case done <- struct{}{}:
		default:
		}
	})))

	// 阻塞到当前片段自然播放完毕。
	// 这样可以保证函数返回时, 该次播放对应的生命周期已完整结束,
	// 避免调用方误以为播放已经完成而继续触发清理或下一步依赖逻辑。
	<-done
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

// splitRouteAudioVolumeNormalProcessing 在“分离路由”模式下叠加键盘/鼠标独立音量。
// 设计原则：
//   - 仅当 SourceMode 为 route-split 时生效
//   - 与主页面全局音量叠加（先全局后分离）
//   - 默认值为 0（不增不减），上限为 0，避免放大超过原始音量
func splitRouteAudioVolumeNormalProcessing(audioVolume *effects.Volume, keycode string) *effects.Volume {
	if audioVolume == nil {
		return audioVolume
	}
	if strings.TrimSpace(keycode) == "" {
		return audioVolume
	}

	state := GetPlaybackState()
	if state.SourceMode != SourceModeRouteSplit {
		return audioVolume
	}

	isMouse := strings.HasPrefix(keycode, "-")
	var configKey string
	var silentKey string
	var defaultValue float64
	var defaultSilent bool

	if isMouse {
		configKey = "main_home.split_audio_volume_processing.mouse.volume_normal"
		defaultValue = config.Main_home___split_audio_volume_processing___mouse_volume_normal
		silentKey = "main_home.split_audio_volume_processing.mouse.volume_silent"
		defaultSilent = config.Main_home___split_audio_volume_processing___mouse_volume_silent
	} else {
		configKey = "main_home.split_audio_volume_processing.keyboard.volume_normal"
		defaultValue = config.Main_home___split_audio_volume_processing___keyboard_volume_normal
		silentKey = "main_home.split_audio_volume_processing.keyboard.volume_silent"
		defaultSilent = config.Main_home___split_audio_volume_processing___keyboard_volume_silent
	}

	volumeValue, ok := config.GetValue(configKey).(float64)
	if !ok {
		volumeValue = defaultValue
		go config.SetValue(configKey, defaultValue)
	}

	// 保证不超过 0，避免放大超过原始音量。
	if volumeValue > 0 {
		volumeValue = 0.0
		go config.SetValue(configKey, 0.0)
	}

	silentValue, ok := config.GetValue(silentKey).(bool)
	if !ok {
		silentValue = defaultSilent
		go config.SetValue(silentKey, defaultSilent)
	}

	return &effects.Volume{
		Streamer: audioVolume,
		Base:     1.6,
		Volume:   volumeValue,
		Silent:   silentValue,
	}
}

// pressReleaseAudioVolumeNormalProcessing 在开启“按下/抬起音量单独控制”时，按事件态叠加独立音量。
// 叠加顺序在调用链中固定为：全局 -> 分离设备 -> 按下/抬起。
func pressReleaseAudioVolumeNormalProcessing(audioVolume *effects.Volume, keycode string, keyState string) *effects.Volume {
	if audioVolume == nil {
		return audioVolume
	}

	if keyState != KeyStateDown && keyState != KeyStateUp {
		return audioVolume
	}

	isEnabled, ok := config.GetValue("main_home.press_release_audio_volume_processing.is_enabled").(bool)
	if !ok {
		isEnabled = config.Main_home___press_release_audio_volume_processing___is_enabled
		go config.SetValue("main_home.press_release_audio_volume_processing.is_enabled", isEnabled)
	}

	if !isEnabled {
		return audioVolume
	}

	buildLayer := func(streamer beep.Streamer, configKey string, silentKey string, defaultValue float64, defaultSilent bool) *effects.Volume {
		volumeValue, ok := config.GetValue(configKey).(float64)
		if !ok {
			volumeValue = defaultValue
			go config.SetValue(configKey, defaultValue)
		}

		if volumeValue > 0 {
			volumeValue = 0.0
			go config.SetValue(configKey, 0.0)
		}

		silentValue, ok := config.GetValue(silentKey).(bool)
		if !ok {
			silentValue = defaultSilent
			go config.SetValue(silentKey, defaultSilent)
		}

		return &effects.Volume{
			Streamer: streamer,
			Base:     1.6,
			Volume:   volumeValue,
			Silent:   silentValue,
		}
	}

	var globalConfigKey string
	var globalSilentKey string
	var globalDefaultValue float64
	var globalDefaultSilent bool
	if keyState == KeyStateDown {
		globalConfigKey = "main_home.press_release_audio_volume_processing.global.down.volume_normal"
		globalSilentKey = "main_home.press_release_audio_volume_processing.global.down.volume_silent"
		globalDefaultValue = config.Main_home___press_release_audio_volume_processing___global_down_volume_normal
		globalDefaultSilent = config.Main_home___press_release_audio_volume_processing___global_down_volume_silent
	} else {
		globalConfigKey = "main_home.press_release_audio_volume_processing.global.up.volume_normal"
		globalSilentKey = "main_home.press_release_audio_volume_processing.global.up.volume_silent"
		globalDefaultValue = config.Main_home___press_release_audio_volume_processing___global_up_volume_normal
		globalDefaultSilent = config.Main_home___press_release_audio_volume_processing___global_up_volume_silent
	}

	layered := buildLayer(audioVolume, globalConfigKey, globalSilentKey, globalDefaultValue, globalDefaultSilent)

	isMouse := strings.HasPrefix(keycode, "-")
	if strings.TrimSpace(keycode) == "" {
		return layered
	}

	var deviceConfigKey string
	var deviceSilentKey string
	var deviceDefaultValue float64
	var deviceDefaultSilent bool
	if isMouse {
		if keyState == KeyStateDown {
			deviceConfigKey = "main_home.press_release_audio_volume_processing.split.mouse.down.volume_normal"
			deviceSilentKey = "main_home.press_release_audio_volume_processing.split.mouse.down.volume_silent"
			deviceDefaultValue = config.Main_home___press_release_audio_volume_processing___split_mouse_down_volume_normal
			deviceDefaultSilent = config.Main_home___press_release_audio_volume_processing___split_mouse_down_volume_silent
		} else {
			deviceConfigKey = "main_home.press_release_audio_volume_processing.split.mouse.up.volume_normal"
			deviceSilentKey = "main_home.press_release_audio_volume_processing.split.mouse.up.volume_silent"
			deviceDefaultValue = config.Main_home___press_release_audio_volume_processing___split_mouse_up_volume_normal
			deviceDefaultSilent = config.Main_home___press_release_audio_volume_processing___split_mouse_up_volume_silent
		}
	} else {
		if keyState == KeyStateDown {
			deviceConfigKey = "main_home.press_release_audio_volume_processing.split.keyboard.down.volume_normal"
			deviceSilentKey = "main_home.press_release_audio_volume_processing.split.keyboard.down.volume_silent"
			deviceDefaultValue = config.Main_home___press_release_audio_volume_processing___split_keyboard_down_volume_normal
			deviceDefaultSilent = config.Main_home___press_release_audio_volume_processing___split_keyboard_down_volume_silent
		} else {
			deviceConfigKey = "main_home.press_release_audio_volume_processing.split.keyboard.up.volume_normal"
			deviceSilentKey = "main_home.press_release_audio_volume_processing.split.keyboard.up.volume_silent"
			deviceDefaultValue = config.Main_home___press_release_audio_volume_processing___split_keyboard_up_volume_normal
			deviceDefaultSilent = config.Main_home___press_release_audio_volume_processing___split_keyboard_up_volume_silent
		}
	}

	return buildLayer(layered, deviceConfigKey, deviceSilentKey, deviceDefaultValue, deviceDefaultSilent)
}

// randomAudioVolumeProcessing 在开启“随机音量”时追加随机衰减层。
// 衰减模型：在当前 volume 体系上减去随机值，deltaVolume = -random(0, maxReduceRatio)。
// 该模型不限制 maxReduceRatio 上限，但始终只会衰减（不放大）。
func randomAudioVolumeProcessing(audioVolume *effects.Volume) *effects.Volume {
	if audioVolume == nil {
		return audioVolume
	}

	isEnabled, ok := config.GetValue("main_home.random_volume_processing.is_enabled").(bool)
	if !ok {
		isEnabled = config.Main_home___random_volume_processing___is_enabled
		go config.SetValue("main_home.random_volume_processing.is_enabled", isEnabled)
	}

	if !isEnabled {
		return audioVolume
	}

	maxReduceRatio, ok := config.GetValue("main_home.random_volume_processing.max_reduce_ratio").(float64)
	if !ok {
		maxReduceRatio = config.Main_home___random_volume_processing___max_reduce_ratio
		go config.SetValue("main_home.random_volume_processing.max_reduce_ratio", maxReduceRatio)
	}

	if maxReduceRatio < 0 {
		maxReduceRatio = 0
		go config.SetValue("main_home.random_volume_processing.max_reduce_ratio", 0.0)
	}

	if maxReduceRatio <= 0 {
		return audioVolume
	}

	randomReduce := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * maxReduceRatio
	deltaVolume := -randomReduce

	return &effects.Volume{
		Streamer: audioVolume,
		Base:     1.6,
		Volume:   deltaVolume,
		Silent:   false,
	}
}

// pressReleaseRandomAudioVolumeProcessing 在开启“按下/抬起随机音量单独控制”时叠加事件态随机层。
// 叠加顺序（调用链）为：全局随机层 -> 按下/抬起随机层（若开启）。
func pressReleaseRandomAudioVolumeProcessing(audioVolume *effects.Volume, keycode string, keyState string) *effects.Volume {
	if audioVolume == nil {
		return audioVolume
	}

	if keyState != KeyStateDown && keyState != KeyStateUp {
		return audioVolume
	}

	isEnabled, ok := config.GetValue("main_home.press_release_random_volume_processing.is_enabled").(bool)
	if !ok {
		isEnabled = config.Main_home___press_release_random_volume_processing___is_enabled
		go config.SetValue("main_home.press_release_random_volume_processing.is_enabled", isEnabled)
	}

	if !isEnabled {
		return audioVolume
	}

	buildRandomLayer := func(streamer beep.Streamer, enabledKey string, ratioKey string, defaultEnabled bool, defaultRatio float64) *effects.Volume {
		nodeEnabled, ok := config.GetValue(enabledKey).(bool)
		if !ok {
			nodeEnabled = defaultEnabled
			go config.SetValue(enabledKey, defaultEnabled)
		}
		if !nodeEnabled {
			return &effects.Volume{Streamer: streamer, Base: 1.6, Volume: 0, Silent: false}
		}

		maxReduceRatio, ok := config.GetValue(ratioKey).(float64)
		if !ok {
			maxReduceRatio = defaultRatio
			go config.SetValue(ratioKey, defaultRatio)
		}
		if maxReduceRatio < 0 {
			maxReduceRatio = 0
			go config.SetValue(ratioKey, 0.0)
		}
		if maxReduceRatio <= 0 {
			return &effects.Volume{Streamer: streamer, Base: 1.6, Volume: 0, Silent: false}
		}

		randomReduce := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * maxReduceRatio
		deltaVolume := -randomReduce

		return &effects.Volume{Streamer: streamer, Base: 1.6, Volume: deltaVolume, Silent: false}
	}

	var globalEnabledKey string
	var globalRatioKey string
	var globalDefaultEnabled bool
	var globalDefaultRatio float64
	if keyState == KeyStateDown {
		globalEnabledKey = "main_home.press_release_random_volume_processing.global.down.is_enabled"
		globalRatioKey = "main_home.press_release_random_volume_processing.global.down.max_reduce_ratio"
		globalDefaultEnabled = config.Main_home___press_release_random_volume_processing___global_down_is_enabled
		globalDefaultRatio = config.Main_home___press_release_random_volume_processing___global_down_max_reduce_ratio
	} else {
		globalEnabledKey = "main_home.press_release_random_volume_processing.global.up.is_enabled"
		globalRatioKey = "main_home.press_release_random_volume_processing.global.up.max_reduce_ratio"
		globalDefaultEnabled = config.Main_home___press_release_random_volume_processing___global_up_is_enabled
		globalDefaultRatio = config.Main_home___press_release_random_volume_processing___global_up_max_reduce_ratio
	}

	layered := buildRandomLayer(audioVolume, globalEnabledKey, globalRatioKey, globalDefaultEnabled, globalDefaultRatio)

	if strings.TrimSpace(keycode) == "" {
		return layered
	}

	isMouse := strings.HasPrefix(keycode, "-")
	var deviceEnabledKey string
	var deviceRatioKey string
	var deviceDefaultEnabled bool
	var deviceDefaultRatio float64

	if isMouse {
		if keyState == KeyStateDown {
			deviceEnabledKey = "main_home.press_release_random_volume_processing.split.mouse.down.is_enabled"
			deviceRatioKey = "main_home.press_release_random_volume_processing.split.mouse.down.max_reduce_ratio"
			deviceDefaultEnabled = config.Main_home___press_release_random_volume_processing___split_mouse_down_is_enabled
			deviceDefaultRatio = config.Main_home___press_release_random_volume_processing___split_mouse_down_max_reduce_ratio
		} else {
			deviceEnabledKey = "main_home.press_release_random_volume_processing.split.mouse.up.is_enabled"
			deviceRatioKey = "main_home.press_release_random_volume_processing.split.mouse.up.max_reduce_ratio"
			deviceDefaultEnabled = config.Main_home___press_release_random_volume_processing___split_mouse_up_is_enabled
			deviceDefaultRatio = config.Main_home___press_release_random_volume_processing___split_mouse_up_max_reduce_ratio
		}
	} else {
		if keyState == KeyStateDown {
			deviceEnabledKey = "main_home.press_release_random_volume_processing.split.keyboard.down.is_enabled"
			deviceRatioKey = "main_home.press_release_random_volume_processing.split.keyboard.down.max_reduce_ratio"
			deviceDefaultEnabled = config.Main_home___press_release_random_volume_processing___split_keyboard_down_is_enabled
			deviceDefaultRatio = config.Main_home___press_release_random_volume_processing___split_keyboard_down_max_reduce_ratio
		} else {
			deviceEnabledKey = "main_home.press_release_random_volume_processing.split.keyboard.up.is_enabled"
			deviceRatioKey = "main_home.press_release_random_volume_processing.split.keyboard.up.max_reduce_ratio"
			deviceDefaultEnabled = config.Main_home___press_release_random_volume_processing___split_keyboard_up_is_enabled
			deviceDefaultRatio = config.Main_home___press_release_random_volume_processing___split_keyboard_up_max_reduce_ratio
		}
	}

	return buildRandomLayer(layered, deviceEnabledKey, deviceRatioKey, deviceDefaultEnabled, deviceDefaultRatio)
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
				}, &Cut{StartMS: 42, EndMS: 60, Volume: 0.6}, keycode, keyState)
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_up.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: -0.6}, keycode, keyState)
			}
		default:
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "test_down.MP3",
				}, &Cut{StartMS: 32, EndMS: 100, Volume: 0}, keycode, keyState)
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "test_up.MP3",
				}, &Cut{StartMS: 28, EndMS: 100, Volume: 0}, keycode, keyState)
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
			}, nil, keycode, keyState)
			return
		}

		if soundEffectType == "sounds" {
			sound_UUID := getValue(configGetter, "key_tone.single."+keycode+"."+keyState+".value")
			if sound_UUID == nil {
				return
			}

			soundParsePlayWith(configGetter, sound_UUID.(string), audioPkgUUID, keycode, keyState)

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
			}, nil, keycode, keyState)
			return
		}

		if soundEffectType == "sounds" {
			sound_UUID := getValue(configGetter, "key_tone.global."+keyState+".value")
			if sound_UUID == nil {
				return
			}

			soundParsePlayWith(configGetter, sound_UUID.(string), audioPkgUUID, keycode, keyState)

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
				}, &Cut{StartMS: 42, EndMS: 60, Volume: 0.6}, keycode, keyState)
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "mouse_test_up.MP3",
				}, &Cut{StartMS: 42, EndMS: 60, Volume: -0.6}, keycode, keyState)
			}
		default:
			if keyState == KeyStateDown {
				PlayKeySound(&AudioFilePath{
					SS: "test_down.MP3",
				}, &Cut{StartMS: 32, EndMS: 100, Volume: 0}, keycode, keyState)
			} else {
				PlayKeySound(&AudioFilePath{
					SS: "test_up.MP3",
				}, &Cut{StartMS: 28, EndMS: 100, Volume: 0}, keycode, keyState)
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
	// 兼容旧调用：默认不区分键盘/鼠标（仅用于非按键上下文）
	soundParsePlayWith(audioPackageConfig.GetValue, sound_UUID, audioPkgUUID, "", "")
}

// audioFileAliasExists 严格校验音频源别名是否存在。
//
// 校验维度：
// 1) `audio_files.<sha256>` 节点存在；
// 2) 节点内 `type` 与引用的 `fileType` 一致；
// 3) 节点内 `name.<nameID>` 真实存在。
//
// 设计目的：
// - 彻底避免“只按 sha256+type 就能播放”的隐式回连；
// - 让运行时播放语义与前端依赖检测语义保持一致（都基于 sha256 + name_id + type 三元组）。
func audioFileAliasExists(get ConfigGetter, sha256 string, nameID string, fileType string) bool {
	if strings.TrimSpace(sha256) == "" || strings.TrimSpace(nameID) == "" || strings.TrimSpace(fileType) == "" {
		return false
	}

	storedType, ok := getValue(get, "audio_files."+sha256+".type").(string)
	if !ok || strings.TrimSpace(storedType) == "" || storedType != fileType {
		return false
	}

	alias := getValue(get, "audio_files."+sha256+".name."+nameID)
	if alias == nil {
		return false
	}

	if aliasText, ok := alias.(string); ok {
		return strings.TrimSpace(aliasText) != ""
	}

	return true
}

func soundParsePlayWith(get ConfigGetter, sound_UUID string, audioPkgUUID string, keycode string, keyState string) {
	sha256, ok := getValue(get, "sounds."+sound_UUID+".source_file_for_sound"+".sha256").(string)
	if !ok {
		logger.Error("message", "error: sha256 value is nil or not a string")
		return
	}

	nameID, ok := getValue(get, "sounds."+sound_UUID+".source_file_for_sound"+".name_id").(string)
	if !ok {
		logger.Error("message", "error: name_id value is nil or not a string")
		return
	}

	fileType, ok := getValue(get, "sounds."+sound_UUID+".source_file_for_sound"+".type").(string)
	if !ok {
		logger.Error("message", "error: file type value is nil or not a string")
		return
	}

	// 关键行为：仅当三元引用（sha256 + name_id + type）真实存在时才允许播放。
	// 这样可确保“删除后重导入同文件”不会自动恢复历史裁剪声音的依赖。
	if !audioFileAliasExists(get, sha256, nameID, fileType) {
		logger.Error("message", "error: source_file_for_sound alias missing",
			"sha256", sha256,
			"name_id", nameID,
			"type", fileType,
		)
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
	}, cut, keycode, keyState)
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
					sha256, shaOK := valueMap["sha256"].(string)
					nameID, idOK := valueMap["name_id"].(string)
					fileType, typeOK := valueMap["type"].(string)
					if !shaOK || !idOK || !typeOK || !audioFileAliasExists(get, sha256, nameID, fileType) {
						logger.Error("message", "error: key_sound single audio_files alias missing",
							"key_sound_uuid", key_sound_UUID,
							"sha256", valueMap["sha256"],
							"name_id", valueMap["name_id"],
							"type", valueMap["type"],
						)
						return
					}
					audio_file_name := sha256 + fileType
					audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
					PlayKeySound(&AudioFilePath{
						Global: audio_file_path,
					}, nil, keycode, keyState)
					return
				}
				if vMap["type"] == "sounds" {
					sound_UUID := vMap["value"].(string)
					soundParsePlayWith(get, sound_UUID, audioPkgUUID, keycode, keyState)
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
				sha256, shaOK := valueMap["sha256"].(string)
				nameID, idOK := valueMap["name_id"].(string)
				fileType, typeOK := valueMap["type"].(string)
				if !shaOK || !idOK || !typeOK || !audioFileAliasExists(get, sha256, nameID, fileType) {
					logger.Error("message", "error: key_sound random audio_files alias missing",
						"key_sound_uuid", key_sound_UUID,
						"sha256", valueMap["sha256"],
						"name_id", valueMap["name_id"],
						"type", valueMap["type"],
					)
					return
				}
				audio_file_name := sha256 + fileType
				audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
				PlayKeySound(&AudioFilePath{
					Global: audio_file_path,
				}, nil, keycode, keyState)
				return
			}
			if vMap["type"] == "sounds" {
				sound_UUID := vMap["value"].(string)
				soundParsePlayWith(get, sound_UUID, audioPkgUUID, keycode, keyState)
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
				sha256, shaOK := valueMap["sha256"].(string)
				nameID, idOK := valueMap["name_id"].(string)
				fileType, typeOK := valueMap["type"].(string)
				if !shaOK || !idOK || !typeOK || !audioFileAliasExists(get, sha256, nameID, fileType) {
					logger.Error("message", "error: key_sound loop audio_files alias missing",
						"key_sound_uuid", key_sound_UUID,
						"sha256", valueMap["sha256"],
						"name_id", valueMap["name_id"],
						"type", valueMap["type"],
					)
					return
				}
				audio_file_name := sha256 + fileType
				audio_file_path := filepath.Join(audioPackageConfig.AudioPackagePath, audioPkgUUID, "audioFiles", audio_file_name)
				PlayKeySound(&AudioFilePath{
					Global: audio_file_path,
				}, nil, keycode, keyState)
				return
			}
			if vMap["type"] == "sounds" {
				sound_UUID := vMap["value"].(string)
				soundParsePlayWith(get, sound_UUID, audioPkgUUID, keycode, keyState)
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
