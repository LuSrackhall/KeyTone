#
# This file is part of the KeyTone project.
#
# Copyright (C) 2024 LuSrackhall
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#


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
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.4.1
	goversioninfo  -skip-versioninfo -platform-specific=true -icon="../frontend/src-electron/icons/icon.ico"
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/KeyTone.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: linux # make linux
linux:
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/KeyTone.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: mac # make mac
mac:
	go build  -ldflags '-s -w' -o "../frontend/dist/key_tone_sdk/KeyTone.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk