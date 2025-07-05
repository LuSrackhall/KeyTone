/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package keyEvent

import (
	"KeyTone/keySound"
	"sync"

	hook "github.com/robotn/gohook"
)

// 定义sse相关变量
type Store struct {
	Keycode uint16 `json:"keycode"`
	State   string `json:"state"`
}

var Clients_sse_stores sync.Map
var once_stores sync.Once

func KeyEventListen() {
	evChan := hook.Start()
	defer hook.End()

	keycode_keycodeChan_map := make(map[uint16]chan hook.Event)
	keycode_buttonChan_map := make(map[uint16]chan hook.Event)

	for ev := range evChan {
		// println("keyAll=", ev.String(), "|||||", ev.Keycode)

		// 防止Dik为0的未知按键(前端一般称其为Unidentified)意外触发。
		if ev.Keycode != 0 {
			// if ev.Kind == 3 || ev.Kind == 4 || ev.Kind == 5 {
			if ev.Kind == 4 || ev.Kind == 5 {
				/* Kind
				 *
				 * KeyDown = 3  // 由于goHook的bug, KeyDown事件对象中, 无法判断实际的Keycode(即Keycode始终为0), 因此我们不使用这个事件。转而利用KeyHold事件代替此事件。
				 * KeyHold = 4
				 * KeyUp   = 5
				 *
				 * 鼠标的这些我们不需要,  只关注 button 即可
				 * MouseUp    = 6        // 由于MouseUp存在问题(Hold时间过长时, Up事件会失效)。
				 * MouseHold  = 7        // 不必多说, 这就是Down了。
				 * MouseDown  = 8				 // 因触发特性是抬起触发, 且不存在Hold时间过长时的失效问题, 故可作为Up使用。
				 * MouseMove  = 9        // 鼠标移动时
				 * MouseDrag  = 10       // 鼠标拖动时(如在桌面按左/右键移动, 会出现选框, 此时会触发此事件[稍后还需确认macos上是否会触发])
				 * MouseWheel = 11       // 鼠标滚轮滚动时, 会触发此事件。
				 */
				// println("keyAll=", ev.String(), "|||||", ev.Keycode)
				// if ev.Kind == 3 {
				// 	println("down")
				// 	println(ev.Keycode) // 按下时, 由于goHook的bug, 故无法判断实际的Keycode, 因此我们不使用这个事件。
				// }
				if _, exists := keycode_keycodeChan_map[ev.Keycode]; exists {
					// logger.Debug("此时已经有了处理此按键发音的通道与其专用的goroutine, 因此无需进行任何创建操作, 只需要向其传递最新的事件信号即可")
					keycode_keycodeChan_map[ev.Keycode] <- ev
				} else {
					// logger.Debug("此时还没有处理此按键发音的通道与其专用的goroutine, 因此需进行相关的创建操作, 并在创建后向其传递最新的事件信号")
					// 创建此按键的专属通道channel
					keycode_keycodeChan_map[ev.Keycode] = make(chan hook.Event)
					// 创建此按键专属 按键事件处理 的 goroutine
					go handleKeyEvent(keycode_keycodeChan_map[ev.Keycode])
					// 将本次按键事件传递至相关通道channel
					keycode_keycodeChan_map[ev.Keycode] <- ev

				}
			}
		} else {
			// 鼠标事件的Keycode等于0
			if ev.Kind == 7 || ev.Kind == 8 {

				// println("keyAll=", ev.String(), "|||||", ev.Keycode)
				if _, exists := keycode_buttonChan_map[ev.Button]; exists {
					keycode_buttonChan_map[ev.Button] <- ev
				} else {
					keycode_buttonChan_map[ev.Button] = make(chan hook.Event)
					// 创建此按键专属 按键事件处理 的 goroutine
					go handleKeyEvent(keycode_buttonChan_map[ev.Button])
					// 将本次按键事件传递至相关通道channel
					keycode_buttonChan_map[ev.Button] <- ev
				}

			}
		}

	}
}

func handleKeyEvent(evChan chan hook.Event) {

	var key_down_soundIsRun bool = false

	for ev := range evChan {
		if ev.Kind == 4 {
			// println("hold")
			// println(ev.Keycode) // 按下时, 由于goHook的bug, 无法判断实际的keyCode。 但由于hold的触发实际与down几乎一致, 且可判断实际的keyCode, 因此可使用此事件代替down
			if !key_down_soundIsRun {
				println("")
				println("")
				println("=====down=====")
				println("   ********")
				println("hold | down ======>", "keycode=", ev.Keycode, "  down")
				// println("keyAll=", ev.String())
				println("仅播放 key_down 声音")
				println("   ********")
				println("=====down=====")
				println("")
				// go keySound.PlayKeySound("test_down.MP3")
				// go keySound.PlayKeySound(&keySound.AudioFilePath{
				// 	SS: "test_down.MP3",
				// }, nil)
				// go keySound.KeyDownSoundPlay()

				go keySound.KeySoundHandler(keySound.KeyStateDown, ev.Keycode)
				key_down_soundIsRun = true
				go sseBroadcast(&Clients_sse_stores, &Store{
					Keycode: ev.Keycode,
					State:   keySound.KeyStateDown,
				})
			}
		}

		if ev.Kind == 5 {

			println("")
			println("")
			println("======up======")
			println("   ********")
			println("up          ======>", "keycode=", ev.Keycode, "  up")
			// println("keyAll=", ev.String())
			println("仅播放 key_up 声音")
			println("   ********")
			println("======up======")
			println("")
			// go keySound.PlayKeySound("test_up.MP3")
			// TODO: 第一个参数何时为nil, 由配置决定。(比如可以设置个bool值, 代表是否关闭此音频, 如果为true, 则为nil。)
			// TODO: 第二个参数何时为nil, 由配置决定。(比如当配置中, 未对开始, 结束时间做任何设置, 则默认为nil的全量播放)
			// TODO: 第一个参数 与 第二个参数组合的更多逻辑, 都需要逐渐适配。比如
			//       * 配置中开启优先级时, 则按照优先级指定音频, 与默认音频, 否则将放弃相关逻辑
			//         * 当开启优先级时, 第一个参数的全局配置, 应该是和相应的cut是强关联的, 否则无意义, 因此PlayKeySound的逻辑还要改。
			// go keySound.PlayKeySound(&keySound.AudioFilePath{
			// 	SS: "test_up.MP3",
			// }, nil) // 注意, 若第二个参数为nil, 则不论多长的音频, 都会全量播放
			// go keySound.KeyUpSoundPlay()
			go keySound.KeySoundHandler(keySound.KeyStateUp, ev.Keycode)

			key_down_soundIsRun = false

			go sseBroadcast(&Clients_sse_stores, &Store{
				Keycode: ev.Keycode,
				State:   keySound.KeyStateUp,
			})
		}

		// var mouse_key_down_soundIsRun bool = false // mouse的hold不存在持续触发的问题, 因此不需要此变量

		if ev.Kind == 7 {
			// println("buttonDown=", ev.Button)

			println("")
			println("")
			println("=====down=====")
			println("   ********")
			println("hold | down ======>", "mouse_button=", ev.Button, "  down")
			println("仅播放 mouse_key_down 声音")
			println("   ********")
			println("=====down=====")
			println("")

			go keySound.KeySoundHandler(keySound.KeyStateMouseDown, ev.Button)
			// mouse_key_down_soundIsRun = true
			go sseBroadcast(&Clients_sse_stores, &Store{
				Keycode: ev.Button,
				State:   keySound.KeyStateMouseDown,
			})
		}

		if ev.Kind == 8 {
			// println("buttonUp=", ev.Button)

			println("")
			println("")
			println("======up======")
			println("   ********")
			println("up          ======>", "mouse_button=", ev.Button, "  up")
			println("仅播放 mouse_key_up 声音")
			println("   ********")
			println("======up======")
			println("")

			go keySound.KeySoundHandler(keySound.KeyStateMouseUp, ev.Button)

			go sseBroadcast(&Clients_sse_stores, &Store{
				Keycode: ev.Button,
				State:   keySound.KeyStateMouseUp,
			})
		}
	}
}

func sseBroadcast(Clients_sse_stores *sync.Map, store *Store) {
	Clients_sse_stores.Range(func(key, value interface{}) bool {
		clientChan := key.(chan *Store)
		serverChan := value.(chan bool)
		select {
		case clientChan <- store:
			return true
		case <-serverChan:
			once_stores.Do(func() {
				close(serverChan)
			})
			return true
		}
	})
}
