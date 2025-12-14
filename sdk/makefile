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


# 允许通过环境变量覆盖/追加 ldflags，便于注入混淆后的密钥
# ?= 语法说明：条件赋值。
# 仅当变量 BASE_LDFLAGS 未被定义时，才将其赋值为 -s -w。
# 如果用户在执行 make 时已经设置了该变量（如 make BASE_LDFLAGS="..."），则使用用户设置的值。
# -s: 省略符号表和调试信息，减小二进制体积
# -w: 省略 DWARF 符号表，减小二进制体积
BASE_LDFLAGS ?= -s -w

# EXTRA_LDFLAGS 用于接收外部注入的额外链接参数（如密钥注入）
# 默认为空。
# make 会自动继承 shell 中的环境变量，因此如果在执行 make 前 export EXTRA_LDFLAGS="..."，
# 这里就会自动获取到该值，无需手动传递。
EXTRA_LDFLAGS ?=


.PHONY: win # make win
# win: spa
win:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.4.1
	goversioninfo -platform-specific=true -icon="../frontend/src-electron/icons/icon.ico"  -manifest="resource\goversioninfo.exe.manifest"
	go build  -ldflags '$(BASE_LDFLAGS) $(EXTRA_LDFLAGS)' -o "../frontend/dist/key_tone_sdk/KeyTone.exe"
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: linux # make linux
linux:
	go build  -ldflags '$(BASE_LDFLAGS) $(EXTRA_LDFLAGS)' -o "../frontend/dist/key_tone_sdk/KeyTone.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk

.PHONY: mac # make mac
mac:
	go build  -ldflags '$(BASE_LDFLAGS) $(EXTRA_LDFLAGS)' -o "../frontend/dist/key_tone_sdk/KeyTone.exe" main.go
	cd ../frontend && quasar build -m electron && rm -rf dist/key_tone_sdk