# .PHONY: build-spa # make build-spa
# build-spa:
# 	cd ../frontend && quasar build
# 	cd ../sdk && rm -rf key_tone_sdk/spa && cp -r ../frontend/dist/spa key_tone_sdk/spa
.PHONY: build-Windows # make build-Windows
# build-Windows: build-spa
build-Windows:
	GOOS=windows GOARCH=amd64 go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/key_tone_sdk.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: all # make all
all: build-Windows