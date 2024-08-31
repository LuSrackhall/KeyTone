package main

import (
	"KeyTone/config"
	"KeyTone/keySound"
	"KeyTone/logger"
	"KeyTone/server"
	"flag"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	hook "github.com/robotn/gohook"
)

func init() {

	// 设置环境变量
	err := godotenv.Load()
	// 在 InitLogger 之前, 使用slog的log信息(与我们的logger模块无关), 此类信息仅在终端上看即可, 输出到日志中也无意义。
	if err != nil {
		// 没必要因为这个err退出程序, .env文件在本项目中, 主要作为开发文件使用。(后面真要上配置文件的话, 也是使用.json格式的)
		slog.Warn("无法加载.env文件", "err", err)
	} else {
		slog.Info(".env文件已被正确加载", "SDK_MODE", os.Getenv("SDK_MODE"))
	}

	// 定义配置文件路径的命令行参数
	var configPath string
	// 定义日志文件路径的命令行参数
	var logPathAndName string

	// 获取命令行参数中的传入值
	{

		// 如果路径不存在, 则使用当前目录作为路径
		// * 第一个参数是指向一个字符串变量的指针，用于存储解析后的值。
		// * 第二个参数是命令行参数的名称（在命令行中使用）。  用户在使用时 go run main.go -configPath=./path
		// * 第三个参数是默认值（如果用户没有提供这个参数，则使用默认值）。
		// * 第四个参数是这个参数的描述（帮助信息）。
		flag.StringVar(&configPath, "configPath", ".", "Path to the config file")
		flag.StringVar(&logPathAndName, "logPathAndName", "./log.jsonl", "Path and name to the log file")

		// 解析命令行参数
		flag.Parse()

		// 使用命令行参数
		// ...

	}

	// 初始化模块
	{

		// 初始化日志模块(并顺便初始化gin的MODE), 主要是为了输出到日志中, 便于在用户使用过程中记录bug数据。
		{
			logger.InitLogger(logPathAndName)

			// 设置日志级别(此处主要用于开发过程中, 自己可随时进行调整的级别设置)
			logger.ProgramLevel.Set(slog.LevelDebug)
			// logger.ProgramLevel.Set(slog.LevelInfo)
			// logger.ProgramLevel.Set(slog.LevelWarn)
			// logger.ProgramLevel.Set(slog.LevelError)

			if os.Getenv("SDK_MODE") != "debug" {
				// 设置log库, 在正式release中的默认级别
				logger.ProgramLevel.Set(slog.LevelInfo)

				// 设置 gin 框架, 在正式release中的 MODE 为 "release"
				gin.SetMode(gin.ReleaseMode)
			}

			logger.Info("日志模块已开始正常运行, Getenv值已获取。 ", "SDK_MODE", os.Getenv("SDK_MODE"), "GIN_MODE", os.Getenv("GIN_MODE"))
		}

		// 初始化配置模块
		{
			// 检查指定的路径是否存在
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// 如果路径不存在，创建路径
				err := os.MkdirAll(configPath, os.ModePerm)
				if err != nil {
					logger.Error("配置文件路径创建时出错。", "err", err.Error())
				} else {
					logger.Info("配置文件路径创建成功。", "你的配置文件路径为", configPath)
				}
			} else if err != nil {
				logger.Error("检查配置文件路径时出错。", "err", err.Error())
			} else {
				logger.Info("配置文件路径已存在且无异常。", "你的配置文件路径为", configPath)
			}
			config.ConfigRun(configPath)
		}

	}

}

func main() {
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

func handleKeyEvent(evChan chan hook.Event) {

	var key_down_soundIsRun bool = false

	for ev := range evChan {
		if ev.Kind == 4 {
			println("hold")
			println(ev.Keycode) // 按下时, 由于goHook的bug, 无法判断实际的keyCode。 但由于hold的触发实际与down几乎一致, 且可判断实际的keyCode, 因此可使用此事件代替down
			if !key_down_soundIsRun {
				println("仅播放 key_down 声音")
				// go keySound.PlayKeySound("test_down.wav", nil)
				go keySound.PlayKeySound(&keySound.AudioFilePath{
					SS: "sound.ogg",
				}, &keySound.Cut{
					StartMS: 28393,
					EndMS:   28593,
					// 其它两个字段为设置, 则为"0值"。 bool的"0值"为 false。 int的"0值"为 0
				})
				// go keySound.KeyDownSoundPlay()

				key_down_soundIsRun = true
			}
		}

		if ev.Kind == 5 {

			println("up")
			println(ev.Keycode)

			println("仅播放 key_up 声音")
			// go keySound.PlayKeySound("test_up.wav", nil)
			// TODO: 第一个参数何时为nil, 由配置决定。(比如可以设置个bool值, 代表是否关闭此音频, 如果为true, 则为nil。)
			// TODO: 第二个参数何时为nil, 由配置决定。(比如当配置中, 未对开始, 结束时间做任何设置, 则默认为nil的全量播放)
			// TODO: 第一个参数 与 第二个参数组合的更多逻辑, 都需要逐渐适配。比如
			//       * 配置中开启优先级时, 则按照优先级指定音频, 与默认音频, 否则将放弃相关逻辑
			//         * 当开启优先级时, 第一个参数的全局配置, 应该是和相应的cut是强关联的, 否则无意义, 因此PlayKeySound的逻辑还要改。
			go keySound.PlayKeySound(nil, &keySound.Cut{
				StartMS: 28393,
				// 其它两个字段为设置, 则为"0值"。 bool的"0值"为 false。 int的"0值"为 0
				EndMS: 28593, // 当 EndMS 小于或等于 StartMS  时, 不会播放任何声音
			}) // 注意, 若第二个参数为nil, 则不论多长的音频, 都会全量播放
			// go keySound.KeyUpSoundPlay()

			key_down_soundIsRun = false

		}
	}
}
