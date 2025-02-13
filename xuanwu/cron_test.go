package xuanwu

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/robfig/cron/v3"
)

func Test_Clac_refrsh(t *testing.T) {
	fmt.Println(666)
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse("46 16 * * *"); err != nil {
		// 表达式无效
		fmt.Println("错误")
	} else {
		fmt.Println("成功")
	}
}

/*
测试内存中创建文件,并以shell脚本执行,解决执行脚本内容
*/
func Test_Clac_1(t *testing.T) {
	// 构建一个脚本
	var script bytes.Buffer
	script.WriteString("#!/bin/bash\n")
	script.WriteString(`should_run=true
	counter=0
	while [ $counter -lt 5 ] && [ $should_run = true ]; do
	   echo "当前时间:$(date)"
	   counter=$((counter+1))
	   sleep 1
	done
	echo "脚本执行结束"`)

	// 创建一个bash命令
	cmd := exec.Command("bash", "-s")

	// 将脚本内容传递给标准输入
	cmd.Stdin = &script

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(err.Error())
		return
	}

	// 打印输出
	print(string(output))
}

/** 执行shell内容 **/
func Test_execShellContent(t *testing.T) {
	// _, Writer := mylog.LogInit("/main.log")
	// log.SetOutput(Writer) // 设置默认logger
	// outputFile, _ := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	scriptContent := "should_run=true\ncounter=0\n\nwhile [ $counter -lt 5 ] && [ $should_run = true ]; do\n   echo \"当前时间:$(date)\"\n   counter=$((counter+1))\n   sleep 1\ndone\n\necho \"脚本执行结束\""

	// 在Windows上执行命令需要使用cmd
	shell := "bash"
	if runtime.GOOS == "windows" {
		shell = "cmd"
	}

	// 生成执行命令
	cmd := exec.Command(shell, "-c", scriptContent)
	// cmd.Stderr = outputFile
	// cmd.Stderr = outputFile
	// 执行脚本
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	// err := cmd.Run()
	fmt.Println(err)
	// 打印输出
	log.Println(string(output))
}
func Test_execShellContent1(t *testing.T) {
	// 生成执行命令
	cmd := exec.Command("bash", "-c", "dir")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	// err := cmd.Run()
	fmt.Println(err)
	// 打印输出
	log.Println(string(output))
}

// 执行shell脚本 脚本路径
func Test_execShell(t *testing.T) {

	scriptPath := "./shell.sh"
	var shell string

	if runtime.GOOS == "windows" {
		shell = "cmd"
		scriptPath = strings.ReplaceAll(scriptPath, "/", "\\")
	} else {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, scriptPath)
	// stdout, _ := cmd.StdoutPipe()
	// // output, err := cmd.CombinedOutput()

	// output1 := make([]byte, 1024)
	// for {
	// 	n, _ := stdout.Read(output1)

	// 	if n > 0 {
	// 		log.Print(string(output1))
	// 	}

	// 	if err := cmd.Wait(); err != nil {
	// 		// Command finished
	// 		break
	// 	}
	// }
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
	}

	// 开始执行命令
	if err := cmd.Start(); err != nil {
	}

	// 异步读取输出并写入日志
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil {
				break
			}
			fmt.Println(line)
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
	}

}
