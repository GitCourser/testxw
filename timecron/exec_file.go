package timecron

import (
	"bufio"
	"log"
	"os/exec"

	"golang.org/x/text/encoding/simplifiedchinese"
)

/*
	执行shell脚本文件

* 	脚本文件可以异步执行
*   脚本路径和脚本类型 shell node python
*/
func taskShellFile(filePath string, file_type string, log *log.Logger, callback func()) {
	// 执行shell脚本
	res, err := execShellFile(filePath, file_type, log)

	if err != nil {
		log.Printf("执行脚本%s失败: %v", filePath, err)
	} else {
		log.Printf("执行脚本%s成功,%s", filePath, res)

		// 执行回调函数
		if callback != nil {
			callback()
		}
	}
}

// 执行shell脚本 脚本路径
func execShellFile(scriptPath string, file_type string, log *log.Logger) (string, error) {

	var shell string

	if file_type == "2" { // shell
		shell = "/bin/sh"
	} else if file_type == "5" { //python
		shell = "python"
		// scriptPath = strings.ReplaceAll(scriptPath, "/", "\\")
	} else if file_type == "6" { //python
		shell = "node"
	} else {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell, scriptPath)
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	// 开始执行命令
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// 异步读取输出并写入日志
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil {
				break
			}
			// GBK转utf8
			res, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(line))
			if err != nil {
				log.Printf("编码 %+v 失败, err:%s\n", string(line), err)
			}
			log.Println(string(res))
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return "", err
	}

	// errorString := fmt.Sprintf("%s: %s", err, string(output))

	// return string(output), err
	return "", err
}
