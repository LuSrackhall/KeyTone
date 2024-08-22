# .PHONY: spa # make spa
# spa:
# 	cd ../frontend && quasar build
# 	cd ../sdk && rm -rf key_tone_sdk/spa && cp -r ../frontend/dist/spa key_tone_sdk/spa

# 由于我将使项目依赖与github actions 的运行器, 因此无需再对指定sdk的交叉编译进行指定。 
# > 即无需在go build之前指定交叉编译相关的系统名和cpu架构名了。
# > * 即 `GOOS=windows GOARCH=amd64 go build ......` -> `go build ......`
# 因为接下来, 每个平台以及每个cpu架构的指令集依赖, 都将是原生构建。
# 也正因为如此, 我们才能拥有如此简洁的makefile

.PHONY: win # make win
# win: spa
win:
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/key_tone_sdk.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: linux # make linux
linux:
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/key_tone_sdk.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: mac # make mac
mac:
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/key_tone_sdk.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk