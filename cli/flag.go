package cli

import (
	"flag"
	"fmt"
)

func InitFlag() {

	// 定义命令行参数
	wordPtr := flag.String("list", "foo", "输出定时任务列表")

	numbPtr := flag.Int("stop", 0, "停止所有任务")
	boolPtr := flag.String("run", "", "运行定时任务")

	// 解析命令行参数
	flag.Parse()

	// 获取参数值
	fmt.Println("list:", *wordPtr)
	fmt.Println("stop:", *numbPtr)
	fmt.Println("run:", *boolPtr)
}
