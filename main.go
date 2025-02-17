package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"syscall"
	"time"

	"xuanwu/config"
	serve "xuanwu/gin"
	xwlog "xuanwu/log"
	"xuanwu/xuanwu"
)

func init() {
	if runtime.GOOS == "linux" { //windows上设置时区会报错,不设置也会正常显示,linux日志时间会差8小时
		TIME_LOCATION, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			log.Printf("time时区设置失败")
			panic(err)
		}
		time.Local = TIME_LOCATION
	}
}

func main() {
	// 添加命令行参数
	hideWindow := flag.Bool("hide", false, "在Windows平台下隐藏命令提示符窗口")
	flag.Parse()

	// 如果是Windows平台且启用了hide参数，则隐藏窗口
	if runtime.GOOS == "windows" && *hideWindow {
		if kernel32 := syscall.NewLazyDLL("kernel32.dll"); kernel32 != nil {
			if proc := kernel32.NewProc("GetConsoleWindow"); proc != nil {
				hwnd, _, _ := proc.Call()
				if hwnd != 0 {
					user32 := syscall.NewLazyDLL("user32.dll")
					if user32 != nil {
						proc := user32.NewProc("ShowWindow")
						proc.Call(hwnd, 0) // 0 表示隐藏窗口
					}
				}
			}
		}
	}

	//初始化日志文件
	_, Writer := xwlog.LogInit("main.log")
	log.SetOutput(Writer) // 设置默认logger

	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	fmt.Println(time.Now())
	log.Println("系统main启动")

	//初始化web服务 传递端口
	go serve.InitApi(cfg, nil)
	//初始化定时任务
	xuanwu.CronInit(cfg)
}
