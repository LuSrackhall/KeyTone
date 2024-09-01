package keySound

import (
	"KeyTone/config"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

//go:embed sounds
var sounds embed.FS

type AudioFilePath struct {
	SS     string // 优先级最低
	Global string // 优先级仅次于Part
	Part   string // 优先级最高
}

type Cut struct {
	StartMS int
	EndMS   int // 当 EndMS 小于或等于 StartMS  时, 不会播放任何声音
}

// 键音播放器
//
// Parameters:
//   - audioFilePath - 指定音频文件路径的结构体, 为nil代表不播放任何音频。
//   - cut - 裁剪键音的必要结构体, 为nil代表不裁剪。
//
// Returns:
//   - void
func PlayKeySound(audioFilePath *AudioFilePath, cut *Cut) {
	// speaker.Clear()
	// defer speaker.Clear()
	if audioFilePath == nil {
		return
	}

	var audioFile fs.File
	var err error  // 注意, 这里一定要同时带上err。 否则在if else 内部, 和已声明的audioFile一起取返回值而临时创建的err, 会造成已声明的audioFile被重新声明并定义, 从而发生作用域问题。
	var ext string // 用于判断音频类型
	if audioFilePath.Part != "" {
		audioFile, err = os.Open(audioFilePath.Part)
		if err != nil {
			panic(err)
		}
		defer audioFile.Close()
		ext = strings.ToLower(filepath.Ext(audioFilePath.Part))
	} else if audioFilePath.Global != "" {
		audioFile, err = os.Open(audioFilePath.Global)
		if err != nil {
			panic(err)
		}
		defer audioFile.Close()
		ext = strings.ToLower(filepath.Ext(audioFilePath.Global))
	} else {
		audioFile, err = sounds.Open("sounds/" + audioFilePath.SS)
		if err != nil {
			panic(err)
		}
		defer audioFile.Close()
		ext = strings.ToLower(filepath.Ext(audioFilePath.SS))
	}

	// 对文件进行解码
	audioStreamer, format, err := decodeAudioFile(audioFile, ext)
	if err != nil {
		panic(err)
	}
	defer audioStreamer.Close()

	// 初始化speaker。
	// 第二个参数的值, 不会对音质产生影响, 它只是缓冲区的大小。
	// > bufferSize参数指定扬声器缓冲区的样本数。更大的缓冲区大小意味着更低的CPU使用率和更可靠的播放。较低的缓冲区大小意味着更好的响应性和更少的延迟。
	// > * 缓冲区越大, cpu压力越小, 播放的整个过程崩溃率也会降低。(个人理解)
	// > * 缓冲区越小, cpu压力越大, 会得到更快的响应性和更少的延时。(个人理解)
	// > 鉴于个人的以上理解, 这个数值对我们KeyTone项目来说, 缓冲区设置的越小越好。
	// > * 但实际测试下来, 缓冲区无论如何设置, 其响应到播放完毕的用时都只有最大20ms作用的波动, 而且绝大部分时候, 波动仅有1ms左右。因此给其一个固定的值即可
	// starTime := time.Now()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/36))

	endSample := -1 // 为保证cut=nil时, 也能正常保留原始工作。(当从配置文件获取的信息达不到构造cut时, cut将不会被构造。cut释放为nil的逻辑不应该在播放器端处理<如start和end都等于0时, cut就应该为nil, 即全量PlayKeySound播放>。)
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
	}

	volume := globalAudioVolumeAmplifyProcessing(audioStreamer)

	volume = globalAudioVolumeNormalProcessing(volume)

	ctrl := &beep.Ctrl{Streamer: volume, Paused: false}

	// 播放音乐
	done := make(chan bool)
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	})))

	// 等待播放完成
	// <-done
	//
	// 等待播放完成
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case <-time.After(time.Second):
	// 		speaker.Lock()
	// 		fmt.Println(format.SampleRate.D(audioStreamer.Position()).Round(time.Second))
	// 		fmt.Println(audioStreamer.Position())
	// 		fmt.Println(format.SampleRate.D(audioStreamer.Len()).Round(time.Second))
	// 		fmt.Println(audioStreamer.Len())
	// 		speaker.Unlock()
	// 	}
	// }
	//
	// 等待播放完成
	re := true
	for re {
		select {
		case <-done:
			re = false
		case <-time.After(10 * time.Millisecond):
			// speaker.Lock()
			pos := audioStreamer.Position()
			fmt.Println("pos", pos)
			fmt.Println("endSample", endSample)
			if pos >= endSample && endSample != -1 {
				// audioStreamer.Close()
				// done <- true
				// return
				ctrl.Paused = true // 目前只能用此一种方式, 在指定时间中止正在播放的音频
				// speaker.Lock()
				// speaker.Clear()
				err := audioStreamer.Close()
				if err != nil {
					fmt.Println(err)
				}
				err = audioFile.Close()
				if err != nil {
					fmt.Println(err)
				}
				re = false
			}
			// speaker.Unlock()
		}
		// fmt.Println("播放用时", time.Since(starTime))
	}
	fmt.Println("退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,退出了,")

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
	initKeyDownSoundBuffer()
	initKeyUpSoundBuffer()
}

var bufferKeyDownSound *beep.Buffer
var bufferKeyUpSound *beep.Buffer

func initKeyDownSoundBuffer() {
	audioFile, err := sounds.Open("sounds/test_down.wav")
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
	audioFile, err := sounds.Open("sounds/test_up.wav")
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
