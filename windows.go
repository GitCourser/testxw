//go:build windows

package main

import (
	"log"
	"syscall"
	// "unsafe"
)

func hideConsoleWindow() {
	log.Println("隐藏控制台流程开始...")
	defer log.Println("隐藏流程结束")
	if !*hideWindow {
		return
	}

	// 隐藏窗口
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd != 0 {
		user32 := syscall.NewLazyDLL("user32.dll")
		showWindow := user32.NewProc("ShowWindow")
		showWindow.Call(hwnd, 0)
	}

	// 释放控制台
	freeConsole := kernel32.NewProc("FreeConsole")
	if ret, _, _ := freeConsole.Call(); ret == 0 {
		log.Println("控制台释放失败（可能已释放）")
	}

	// 修改控制台标题（可选）
	// _, _, _ = syscall.SyscallN(
	// 	kernel32.NewProc("SetConsoleTitleW").Addr,
	// 	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("XuanWu Background Service"))),
	// )
}
