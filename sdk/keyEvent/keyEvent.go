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

	for ev := range evChan {
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
