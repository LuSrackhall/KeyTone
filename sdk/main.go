package main

import (
	"KeyTone/config"
	"KeyTone/keySound"
	"KeyTone/server"
	"flag"
	"fmt"
	"os"

	hook "github.com/robotn/gohook"
)

func main() {

	// 获取命令行参数中的配置文件路径
	{
		// 如果不需要默认参数, 则可直接使用os包来获取参数变量
		// config.ConfigRun(os.Args[1])

		// 如果需要默认参数, 则使用专门的库来解析获取参数变量
		// * 此功能需要使用命令行参数解析库, 最简单的有个官方库flag
		// * 对于复杂的命令行交互, 推荐使用第三方库[cobra](https://github.com/spf13/cobra) 或 [cli](https://github.com/urfave/cli)

		// 定义命令行参数
		var configPath string

		// 如果路径不存在, 则使用当前目录作为路径
		// * 第一个参数是指向一个字符串变量的指针，用于存储解析后的值。
		// * 第二个参数是命令行参数的名称（在命令行中使用）。  用户在使用时 go run main.go -configPath=./path
		// * 第三个参数是默认值（如果用户没有提供这个参数，则使用默认值）。
		// * 第四个参数是这个参数的描述（帮助信息）。
		flag.StringVar(&configPath, "configPath", ".", "Path to the config file")

		// 解析命令行参数
		flag.Parse()

		// 使用命令行参数
		{
			// 检查指定的路径是否存在
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("指定的路径不存在，使用当前目录")
				configPath = "." // 使用当前目录
			}
			config.ConfigRun(configPath)
		}
	}

	go server.ServerRun()
	keyEventListen()
}

func keyEventListen() {
	evChan := hook.Start()
	defer hook.End()

	keycode_keycodeChan_map := make(map[uint16]chan hook.Event)

	for ev := range evChan {
		// if ev.Kind == 3 || ev.Kind == 4 || ev.Kind == 5 {
		if ev.Kind == 4 || ev.Kind == 5 {
			/* Kind
			 *
			 * KeyDown = 3  // 由于goHook的bug, KeyDown事件对象中, 无法判断实际的Keycode(即Keycode始终为0), 因此我们不使用这个事件。转而利用KeyHold事件代替此事件。
			 * KeyHold = 4
			 * KeyUp   = 5
			 *
			 * 鼠标的这些我们不需要,  只关注 button 即可
			 * MouseUp    = 6
			 * MouseHold  = 7
			 * MouseDown  = 8
			 * MouseMove  = 9
			 * MouseDrag  = 10
			 * MouseWheel = 11
			 */
			// if ev.Kind == 3 {
			// 	println("down")
			// 	println(ev.Keycode) // 按下时, 由于goHook的bug, 故无法判断实际的Keycode, 因此我们不使用这个事件。
			// }
			if _, exists := keycode_keycodeChan_map[ev.Keycode]; exists {
				// fmt.Println("此时已经有了处理此按键发音的通道与其专用的goroutine, 因此无需进行任何创建操作, 只需要向其传递最新的事件信号即可")
				keycode_keycodeChan_map[ev.Keycode] <- ev
			} else {
				// fmt.Println("此时还没有处理此按键发音的通道与其专用的goroutine, 因此需进行相关的创建操作, 并在创建后向其传递最新的事件信号")
				// 创建此按键的专属通道channel
				keycode_keycodeChan_map[ev.Keycode] = make(chan hook.Event)
				// 创建此按键专属 按键事件处理 的 goroutine
				go handleKeyEvent(keycode_keycodeChan_map[ev.Keycode])
				// 将本次按键事件传递至相关通道channel
				keycode_keycodeChan_map[ev.Keycode] <- ev

			}
		}
	}
}

func handleKeyEvent(evChan chan hook.Event) {

	var key_down_soundIsRun bool = false

	for ev := range evChan {
		if ev.Kind == 4 {
			println("hold")
			println(ev.Keycode) // 按下时, 由于goHook的bug, 无法判断实际的keyCode。 但由于hold的触发实际与down几乎一致, 且可判断实际的keyCode, 因此可使用此事件代替down
			if !key_down_soundIsRun {
				println("仅播放 key_down 声音")
				go keySound.PlayKeySound("test_down.wav")
				// go keySound.KeyDownSoundPlay()

				key_down_soundIsRun = true
			}
		}

		if ev.Kind == 5 {

			println("up")
			println(ev.Keycode)

			println("仅播放 key_up 声音")
			go keySound.PlayKeySound("test_up.wav")
			// go keySound.KeyUpSoundPlay()

			key_down_soundIsRun = false

		}
	}
}
