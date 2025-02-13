package xuanwu

import (
	"log"
	"os/exec"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
)

/*

****  一次性,非交互式命令直接可以获取到执行结果 ****
** 非异步命令 ***

 */
/*
* 执行命令行命令
*

  - Go的exec包默认是使用 cmd /c 格式在Windows运行命令
    所以就不用判断平台传递cmd参数,直接写命令

    这里只执行 一次性,非交互式命令,直接等待获取执行结果即可
*/
func exec_command(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	// windows使用-c参数会直接打开一个cmd就获取不到输出结果,需要使用cmd /c
	var shell string
	var Type string

	if runtime.GOOS == "windows" {
		shell = "cmd"
		Type = "/c"
	} else {
		shell = "/bin/sh"
		Type = "-c"
	}
	cmd := exec.Command(shell, Type, s)
	output, err := cmd.CombinedOutput() // 这个会等待一次获取所有的cmd.Stdout,cmd.Stderr,不区分
	if err != nil {
		return "", err
	}
	// GBK转utf8
	res, err := simplifiedchinese.GBK.NewDecoder().Bytes(output)
	if err != nil {
		log.Printf("编码 %+v 失败, err:%s\n", string(output), err)
	}
	return string(res), err
}

/*
执行命令行任务
*/
func taskCommand(command string, log *log.Logger) {
	res, err := exec_command(command)
	if err != nil {
		log.Println("\n", err, "命令行执行失败", "\n")
		return
	}
	log.Println("\n", res, "命令执行成功", "\n")
}
