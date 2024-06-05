package main

import (
	"KeyTone/keySound"

	hook "github.com/robotn/gohook"
)

func main() {

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

				key_down_soundIsRun = true
			}
		}

		if ev.Kind == 5 {

			println("up")
			println(ev.Keycode)

			println("仅播放 key_up 声音")
			go keySound.PlayKeySound("test_up.wav")

			key_down_soundIsRun = false

		}
	}
}
