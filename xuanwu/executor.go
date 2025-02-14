package xuanwu

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 处理工作目录路径
func HandleWorkDir(workDir string) string {
	if workDir == "" {
		return ""
	}
	
	// Windows系统检查盘符
	if runtime.GOOS == "windows" {
		if len(workDir) >= 2 && workDir[1] == ':' {
			return workDir
		}
	} else {
		// Linux/Unix系统检查根目录
		if strings.HasPrefix(workDir, "/") {
			return workDir
		}
	}
	
	// 相对路径处理
	exePath, err := os.Executable()
	if err != nil {
		return workDir
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "data", workDir)
}

// 执行任务命令
func ExecTask(command string, workDir string, log *log.Logger) error {
	fmt.Printf("准备执行命令: %s, 工作目录: %s\n", command, workDir)
	// 处理工作目录
	workDir = HandleWorkDir(workDir)
	fmt.Printf("处理后的工作目录: %s\n", workDir)
	
	// 创建命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	// 设置工作目录
	if workDir != "" {
		cmd.Dir = workDir
		fmt.Printf("设置命令工作目录: %s\n", workDir)
	}
	
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("获取标准输出管道失败: %v\n", err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("获取标准错误管道失败: %v\n", err)
		return err
	}
	
	// 开始执行命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("启动命令失败: %v\n", err)
		return err
	}
	
	// 异步读取标准输出
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("读取输出错误: %v\n", err)
					log.Printf("读取输出错误: %v\n", err)
				}
				break
			}
			fmt.Printf("命令输出: %s", line)
			log.Print(line)
		}
	}()
	
	// 异步读取标准错误
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("读取错误输出错误: %v\n", err)
					log.Printf("读取错误输出错误: %v\n", err)
				}
				break
			}
			fmt.Printf("命令错误输出: %s", line)
			log.Print(line)
		}
	}()
	
	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("命令执行失败: %v\n", err)
	} else {
		fmt.Println("命令执行完成")
	}
	return err
}
