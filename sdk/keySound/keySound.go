package keySound

import (
	"KeyTone/config"
	"embed"
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

	volume := globalAudioVolumeProcessing(audioStreamer)

	// 播放音乐
	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	// 等待播放完成
	<-done

}

func globalAudioVolumeProcessing(audioStreamer beep.Streamer) *effects.Volume {
	var audio_volume_processing_volume_amplify = config.GetValue("audio_volume_processing.volume_amplify")
	volumeAmplify := &effects.Volume{
		Streamer: audioStreamer,
		Base:     1.6,
		Volume:   audio_volume_processing_volume_amplify.(float64),
		Silent:   false,
	}

	var audio_volume_processing_volume_normal = config.GetValue("audio_volume_processing.volume_normal")
	volumeNormal := &effects.Volume{
		Streamer: volumeAmplify,
		Base:     1.6,
		Volume:   audio_volume_processing_volume_normal.(float64),
		Silent:   false,
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
