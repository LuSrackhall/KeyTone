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
/*
printconfig - KeyTone 专辑配置解密查看工具


密钥说明：
  解密使用 FixedSecret 派生 AES key（SHA256(secret + last6(sha1(albumUUID)))）。
  FixedSecret 支持构建时注入（-ldflags -X KeyTone/audioPackage/enc.FixedSecret=...）。
  - 开源构建：使用默认 FixedSecret，只能解密开源构建产物
  - 私有构建：注入私有 FixedSecret 后，可解密对应私有构建产物

构建方式：
	假设在 sdk/ 目录下, 执行该命令, 会在 sdk/ 目录下生成可执行文件 printconfig, 此cli工具可用于解密开源构建产物, 方便开发调试。
	> 目前已将默认产物添加到.gitignore中避免误提交。

  # 开源构建（使用默认密钥）
  go build ./audioPackage/cmd/printconfig

  # 私有构建（需先加载注入参数）
  source ./setup_build_env.sh
  go build -ldflags "$EXTRA_LDFLAGS" ./audioPackage/cmd/printconfig

直接运行（不落地构建产物）：

	# 开源默认密钥（仅能解密开源构建产物）
	go run ./audioPackage/cmd/printconfig --path /path/to/album/uuid

	# 私有注入密钥（可解密对应私有构建产物）
	source ./setup_build_env.sh
	go run -ldflags "$EXTRA_LDFLAGS" ./audioPackage/cmd/printconfig --path /path/to/album/uuid

兼容性说明：
	- FixedSecret 的覆盖发生在编译/链接阶段（-ldflags -X），不是运行时参数。
	- 因此：同一次运行（同一个 go run/go build 产物）无法同时兼容两套密钥。
		若要分别查看开源/私有产物，请使用不同的 ldflags 分别运行/构建。

用途：
  解密并打印键音专辑的 package.json 配置文件内容。
  当用户选择"需要签名"导出专辑时，配置文件会被 AES-GCM 加密，
  此工具用于在终端中查看加密后的真实配置内容，便于调试。

使用方法：
  printconfig --path <albumDir> [--raw]

参数：
  --path    专辑目录路径（包含 package.json 或 stub + core 文件）
  --raw     输出原始密文 hex（不解密）

示例：
  # 解密并查看配置
  printconfig --path /path/to/album/uuid

  # 仅输出密文 hex（不解密）
  printconfig --path /path/to/album/uuid --raw
*/
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apconfig "KeyTone/audioPackage/config"
	apenc "KeyTone/audioPackage/enc"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: printconfig --path <albumDir> [--raw]\n")
	os.Exit(2)
}

func main() {
	var albumPath string
	var raw bool
	flag.StringVar(&albumPath, "path", "", "album directory path containing package.json")
	flag.BoolVar(&raw, "raw", false, "print raw cipher hex if encrypted")
	flag.Parse()

	if albumPath == "" {
		usage()
	}

	stubInfo, pkgRaw, err := apconfig.ReadCoreStubInfo(albumPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read stub error:", err)
		os.Exit(1)
	}

	if stubInfo != nil {
		corePath := filepath.Join(albumPath, stubInfo.Core)
		cipherBytes, err := os.ReadFile(corePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read core error:", err)
			os.Exit(1)
		}
		if raw {
			fmt.Println(hex.EncodeToString(cipherBytes))
			return
		}
		albumUUID := filepath.Base(albumPath)
		plain, err := apenc.DecryptConfigBytes(cipherBytes, albumUUID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "decrypt error:", err)
			os.Exit(1)
		}
		fmt.Println(plain)
		return
	}

	if pkgRaw == nil {
		pkg := filepath.Join(albumPath, "package.json")
		pkgRaw, err = os.ReadFile(pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			os.Exit(1)
		}
	}

	if apenc.IsLikelyHexCipher(pkgRaw) {
		if raw {
			fmt.Println(strings.TrimSpace(string(pkgRaw)))
			return
		}
		albumUUID := filepath.Base(albumPath)
		plain, err := apenc.DecryptConfigHex(strings.TrimSpace(string(pkgRaw)), albumUUID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "decrypt error:", err)
			os.Exit(1)
		}
		fmt.Println(plain)
		return
	}

	fmt.Println(string(pkgRaw))
}
