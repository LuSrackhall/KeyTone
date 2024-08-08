package keySound

import (
	"KeyTone/config"
	"embed"
	"fmt"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

//go:embed sounds
var sounds embed.FS

func PlayKeySound(ss string) {
	audioFile, err := sounds.Open("sounds/" + ss)
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

	// 初始化speaker。(可更改第二个参数的值(越大, 音质越好, 响应越慢。 越小, 音质越差, 响应速度越快。))
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/36))

	volume := audioVolumeProcessing(audioStreamer)

	// 播放音乐
	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	// 等待播放完成
	<-done

}

func audioVolumeProcessing(audioStreamer beep.Streamer) *effects.Volume {

	// 这里并没重新请求
	var audio_volume_processing map[string]interface{} = config.GetValue("audio_volume_processing").(map[string]interface{})
	fmt.Println("aaaaa", audio_volume_processing) // 这里不应该的

	volumeAmplify := &effects.Volume{
		Streamer: audioStreamer,
		Base:     1.6, // 设置成1.6, 可以使得调整更细腻(越接近1越细腻)的同时, 更加吉利。 (常见的做法是设置成2)。
		Volume:   audio_volume_processing["volume_amplify"].(float64),
		Silent:   false,
	}

	volumeNormal := &effects.Volume{
		Streamer: volumeAmplify,
		// TIPS: 在大多数音频应用中，Base 是固定的，因为它定义了音量调整的行为曲线。如果用户需要调整音量，通常是调整 Volume，而不是 Base。Base 的调整通常是在你希望改变音量增益的对数曲线时。例如，如果你希望音量调整曲线更敏感或更平滑，可以调整 Base。但是，这种需求在实际应用中比较少见。
		Base:   1.6,                                                // 设置成1.6, 可以使得调整更细腻(越接近1越细腻)的同时, 更加吉利。 (常见的做法是设置成2)。
		Volume: audio_volume_processing["volume_normal"].(float64), // nil
		Silent: false,
	}

	return volumeNormal
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
