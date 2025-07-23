module KeyTone

go 1.24.5

require (
	github.com/fsnotify/fsnotify v1.8.0
	github.com/gin-gonic/gin v1.10.0
	github.com/gopxl/beep/v2 v2.1.1
	github.com/robotn/gohook v0.41.0
)

require (
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/ebitengine/oto/v3 v3.3.2 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/gin-contrib/cors v1.7.2
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/viper v1.20.1
	github.com/vcaesar/keycode v0.10.1 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

// 替换官方 viper 为我修改的版本
replace github.com/spf13/viper => github.com/LuSrackhall/viper v1.20.1-custom
