package xuanwu

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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
	// 处理工作目录
	workDir = HandleWorkDir(workDir)
	
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
	}
	
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	
	// 使用WaitGroup等待所有输出读取完成
	var wg sync.WaitGroup
	
	// 开始执行命令
	if err := cmd.Start(); err != nil {
		return err
	}
	
	// 异步读取标准输出
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	
	// 异步读取标准错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	
	// 等待命令执行完成
	err = cmd.Wait()
	
	// 等待所有输出读取完成
	wg.Wait()
	
	return err
}
