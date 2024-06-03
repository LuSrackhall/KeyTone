package main

import (
	hook "github.com/robotn/gohook"
)

func main() {

	keyEventListen()
}

func keyEventListen() {
	evChan := hook.Start()
	defer hook.End()

	var time int64 = 0
	var key_down_soundIsRun bool = false

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

				// (为了正确记录每个独立按键的独立时间, 需要使用切片+ev.Keycode来区分不同的按键才行, 防止并发按键的冲突。今天太晚了, 明天吧!)
				time = 0 // 最终将time归0, 以确保下一次按键触发时, 可以正确的记录时间。

				// (为了正确记录每个独立按键的独立时间, 需要使用切片+ev.Keycode来区分不同的按键才行, 防止并发按键的冲突。今天太晚了, 明天吧!)
				key_down_soundIsRun = false

				// 明天的基本思路是, 利用 每个按键独立 channel + go程 的方式, 处理并发的按键冲突问题。
				// 而此处作为对应触发侧, 只需要通过select将信号正确的传递就好了。
			}

		}
	}
}
