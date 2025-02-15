package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"ktalbum-tools/commands"
	"ktalbum-tools/web"
)

func main() {
	// 定义子命令
	extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
	// 暂时注释掉未使用的命令，直到我们实现它们
	// packCmd := flag.NewFlagSet("pack", flag.ExitOnError)
	// infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	// extract 命令的参数
	extractInput := extractCmd.String("in", "", "输入的 .ktalbum 文件")
	extractOutput := extractCmd.String("out", "", "输出的 .zip 文件 (可选)")
	extractVerbose := extractCmd.Bool("v", false, "显示详细信息")

	// 添加 web 命令
	webCmd := flag.NewFlagSet("web", flag.ExitOnError)
	webPort := webCmd.Int("port", 8080, "Web 服务端口")

	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("使用方法:")
		fmt.Println("  ktalbum-tools extract -in <ktalbum文件> [-out <zip文件>] [-v]")
		fmt.Println("  ktalbum-tools pack -in <zip文件> [-out <ktalbum文件>] [-v]")
		fmt.Println("  ktalbum-tools info -in <ktalbum文件>")
		fmt.Println("  ktalbum-tools web -port <端口>")
		os.Exit(1)
	}

	// 解析子命令
	switch os.Args[1] {
	case "extract":
		extractCmd.Parse(os.Args[2:])
		if *extractInput == "" {
			fmt.Println("请指定输入文件: -in <ktalbum文件>")
			os.Exit(1)
		}
		if *extractOutput == "" {
			*extractOutput = strings.TrimSuffix(*extractInput, ".ktalbum") + ".zip"
		}
		if err := commands.Extract(*extractInput, *extractOutput, *extractVerbose); err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}

	case "web":
		webCmd.Parse(os.Args[2:])
		fmt.Printf("启动 Web 服务在端口 %d...\n", *webPort)
		if err := web.StartServer(*webPort); err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}

	// TODO: 实现 pack 和 info 命令

	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		os.Exit(1)
	}
} 