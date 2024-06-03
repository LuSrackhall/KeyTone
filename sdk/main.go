package main

import (
	hook "github.com/robotn/gohook"
)

func main() {

	low()
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
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
		if ev.Kind == 3 {
			println("down")
			println(ev.Keycode) // 按下时, 由于goHook的bug, 故无法判断实际的keyCode, 因此我们不使用这个事件。
		}
		if ev.Kind == 4 {
			println("hold")
			println(ev.Keycode) // 按下时, 由于goHook的bug, 无法判断实际的keyCode。 但由于hold的触发实际与down几乎一致, 且可判断实际的keyCode, 因此可使用此事件代替down
		}
		if ev.Kind == 5 {
			println("up")
			println(ev.Keycode)
		}
	}
}
