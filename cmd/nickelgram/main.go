package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("无法确定程序目录。")
		return
	}
	root := filepath.Dir(exe)
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "debug" {
		if err := runDebug(args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, "本地测试失败："+err.Error())
			os.Exit(1)
		}
		return
	}
	if len(args) == 0 {
		fmt.Println("NickelGram minimal selection test")
		return
	}
	var runErr error
	switch args[0] {
	case "send-selection":
		if len(args) != 2 {
			fmt.Println("选区参数不完整；未发送。")
			return
		}
		runErr = runAction(root, args[1])
	case "send-latest-highlight":
		if len(args) != 1 {
			fmt.Println("参数不完整；未发送。")
			return
		}
		runErr = runLatestHighlight(root)
	case "export-book":
		if len(args) != 1 {
			fmt.Println("参数不完整；未发送。")
			return
		}
		runErr = runExportBook(root)
	default:
		fmt.Println("未知操作；未发送。")
		return
	}
	if runErr != nil {
		fmt.Println("发送失败：" + runErr.Error())
	}
}
