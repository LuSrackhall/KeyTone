package main

import (
	"encoding/json"
	"fmt"

	hook "github.com/robotn/gohook"
)

func main() {
	add()

	low()
}

func add() {
	// fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	// hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
	// 	fmt.Println("ctrl-shift-q")
	// 	hook.End()
	// })

	fmt.Println("--- Please press w---")
	// hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {

	// 	// fmt.Println("hook(KeyDown): ", e)

	// 	ee, _ := json.Marshal(e)
	// 	fmt.Println("hook(KeyDown): ", string(ee))

	// 	fmt.Println("w")
	// })

	hook.Register(hook.KeyUp, []string{"w"}, func(e hook.Event) {
		// 此处存在的问题https://github.com/robotn/gohook/issues/48

		// fmt.Println("hook(KeyUp): ", e)

		ee, _ := json.Marshal(e)
		fmt.Println("hook(KeyUp): ", string(ee))

		fmt.Println("w")
	})

	s := hook.Start()
	<-hook.Process(s)
}

// Event Holds a system event
//
// If it's a Keyboard event the relevant fields are:
// Mask, Keycode, Rawcode, and Keychar,
// Keychar is probably what you want.   (此处对于键盘事件, 我需要Keycode字段而不是Keychar字段)
//
// If it's a Mouse event the relevant fields are:
// Button, Clicks, X, Y, Amount, Rotation and Direction (此处对于鼠标事件, 我只需要Button字段)
func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		/* 这里打印的ev是 `hook.Event`类型的结构体, 在goHook(即此处的hook)中, 此结构体重写了 String接口方法。
		   因此在通过 fmt.Println打印时, 会按照goHook中的String()方法内的判断逻辑,将变量转化为不同的格式的最终字符串到标准输出。
		*/
		// fmt.Println("hook: ", ev)
		e, _ := json.Marshal(ev)
		fmt.Println("hook: ", string(e))
	}
}
