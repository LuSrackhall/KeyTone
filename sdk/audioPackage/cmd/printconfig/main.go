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
