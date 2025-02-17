//go:build windows

package main

import "syscall"

func hideConsoleWindow() {
    kernel32 := syscall.NewLazyDLL("kernel32.dll")
    if proc := kernel32.NewProc("GetConsoleWindow"); proc != nil {
        hwnd, _, _ := proc.Call()
        if hwnd != 0 {
            user32 := syscall.NewLazyDLL("user32.dll")
            if proc := user32.NewProc("ShowWindow"); proc != nil {
                proc.Call(hwnd, 0)
            }
        }
    }
}