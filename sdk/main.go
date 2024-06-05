package main

import (
	"fmt"

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
		if ev.Kind == 3 || ev.Kind == 4 || ev.Kind == 5 {
			/* Kind
			 *
			 * KeyDown = 3
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
			// 	println(ev.Keycode) // 按下时, 由于goHook的bug, 故无法判断实际的keyCode, 因此我们不使用这个事件。
			// }
			if _, exists := keycode_keycodeChan_map[ev.Keycode]; exists {
				// fmt.Println("此时已经有了处理此按键发音的通道与其专用的goroutine, 因此无需进行任何创建操作, 只需要向其传递最新的事件信号即可")
				keycode_keycodeChan_map[ev.Keycode] <- ev
			} else {
				/*
				 * 你可能会注意到, 在第一次按键按下时, 这里打印了两遍。
				 * 原因是, KeyDown=3时, Keycode的值始终为0, 而我们监听的有3这个情况, 也就是创建了一个并未实际使用的通道和goroutine
				 * 之前也提到过, 这是goHook库的一个bug, 因此我们也是使用 KeyHold , 来代替了 KeyDown。
				 * 这里保留KeyDown只是为了等待goHook库的修复, 多余的这个 Keycode为0的 通道 和 goroutine 对我们的程序性能几乎没有任何影响。
				 * * 虽每次按键都会想 keycode 为0 的通道传递对应事件, 但此事件到达对应的 keycode为0的goroutine中后, 仍会因逻辑而被直接抛弃。
				 * * 因为handleKeyEvent中, 只判断了 Kind = 4 和 Kind = 5的处理逻辑, 因此影响不大。
				 * 当然, 如果认为这会对性能造成影响, 可以直接在父逻辑中取消对 Kind = 3 情况的处理。
				 *
				 */
				fmt.Println("此时还没有处理此按键发音的通道与其专用的goroutine, 因此需进行相关的创建操作, 并在创建后向其传递最新的事件信号")
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

	var time int64 = 0
	var key_down_soundIsRun bool = false

	for ev := range evChan {
		if ev.Kind == 4 {
			println("hold")
			if time == 0 {
				time = ev.When.UnixMilli()
			}
			println(ev.Keycode) // 按下时, 由于goHook的bug, 无法判断实际的keyCode。 但由于hold的触发实际与down几乎一致, 且可判断实际的keyCode, 因此可使用此事件代替down
		}
		if (ev.When.UnixMilli()-time) > 300 && ev.Kind == 4 {
			if !key_down_soundIsRun {
				println("仅播放 key_down 声音")
				key_down_soundIsRun = true
			}
		}
		if ev.Kind == 5 {
			println("up")
			println(ev.Keycode)

			if (ev.When.UnixMilli() - time) > 300 {
				println("仅播放 key_up 声音")
			} else {
				// ev.When - time <= 50
				println("播放 key_down 声音 + key_up 声音")
			}

			time = 0 // 最终将time归0, 以确保下一次按键触发时, 可以正确的记录时间。

			key_down_soundIsRun = false

		}
	}
}
