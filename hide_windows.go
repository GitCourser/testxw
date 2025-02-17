//go:build windows

package main

import (
	"log"
	"syscall"
)

func hideConsoleWindow() {
	log.Println("尝试隐藏控制台窗口...")
	if !*hideWindow {
		return
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	if hwnd, _, _ := getConsoleWindow.Call(); hwnd != 0 {
		user32 := syscall.NewLazyDLL("user32.dll")
		showWindow := user32.NewProc("ShowWindow")
		showWindow.Call(
			hwnd,
			uintptr(0), // SW_HIDE
		)
	}
}