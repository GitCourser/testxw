package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"xuanwu/config"
	serve "xuanwu/gin"
	xwlog "xuanwu/log"
	"xuanwu/xuanwu"
)

// 添加Windows命令行参数
var hideWindow = flag.Bool("hide", false, "在Windows平台下隐藏命令提示符窗口")

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
	flag.Parse()

	// Windows平台特定逻辑
	if runtime.GOOS == "windows" && *hideWindow {
		os.Stdout = Writer
		os.Stderr = Writer
		hideConsoleWindow()
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
