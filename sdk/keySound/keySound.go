package keySound

import (
	"embed"
	"time"

	"github.com/gopxl/beep"
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

	// 播放音乐
	done := make(chan bool)
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))

	// 等待播放完成
	<-done

}
