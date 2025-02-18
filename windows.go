//go:build windows

package main

import (
	// "log"
	"syscall"
)

func hideConsoleWindow() {
	// log.Println("隐藏控制台流程开始...")
	// defer log.Println("隐藏流程结束")
	// if !*hideWindow {
	// 	return
	// }

	// 释放控制台
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	freeConsole := kernel32.NewProc("FreeConsole")
	// ret, _, err := freeConsole.Call()
	freeConsole.Call()

	// 隐藏窗口
	// getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	// if hwnd, _, _ := getConsoleWindow.Call(); hwnd != 0 {
	// 	user32 := syscall.NewLazyDLL("user32.dll")
	// 	showWindow := user32.NewProc("ShowWindow")
	// 	showWindow.Call(hwnd, 0)
	// }
}
