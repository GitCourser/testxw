package timecron

import (
	"bufio"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// 执行shell脚本 脚本路径
func execShell(scriptPath string, log *log.Logger) (string, error) {

	var shell string

	if runtime.GOOS == "windows" {
		shell = "cmd"
		scriptPath = strings.ReplaceAll(scriptPath, "/", "\\")
	} else {
		shell = "/bin/sh"
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
			log.Println(line)
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return "", err
	}

	// errorString := fmt.Sprintf("%s: %s", err, string(output))
	return "", err
}
