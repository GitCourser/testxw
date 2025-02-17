//go:build windows

package main

import (
	"flag"
	"syscall"
)

func init() {
	hideWindow := flag.Bool("hide", false, "在Windows平台下隐藏命令提示符窗口")
	flag.Parse()
	
	if *hideWindow {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		proc := kernel32.NewProc("GetConsoleWindow")
		hwnd, _, _ := proc.Call()
		if hwnd != 0 {
			user32 := syscall.NewLazyDLL("user32.dll")
			proc := user32.NewProc("ShowWindow")
			proc.Call(hwnd, 0)
		}
	}
}