package xwlog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xuanwu/lib/pathutil"
)

// LogConfig 日志配置
type LogConfig struct {
	TaskLogFormat bool // true: 任务日志格式(只在第一行显示时间), false: 标准日志格式
}

// 自定义日志写入器
type taskLogWriter struct {
	file     *os.File
	isHeader bool // 是否已写入头部时间
}

func (w *taskLogWriter) Write(p []byte) (n int, err error) {
	if !w.isHeader {
		// 写入时间头
		timeHeader := time.Now().Format("2006-01-02 15:04:05") + "\n"
		if _, err := w.file.WriteString(timeHeader); err != nil {
			return 0, err
		}
		w.isHeader = true
	}
	return w.file.Write(p)
}

func LogInit(name string) (*log.Logger, *os.File) {
	return LogInitWithConfig(name, &LogConfig{TaskLogFormat: false})
}

func LogInitWithConfig(name string, config *LogConfig) (*log.Logger, *os.File) {
	if name == "" { //没有名称时候,返回空日志
		return log.New(os.Stdout, "", 0), nil
	}
	
	logPath := pathutil.GetLogPath(name)
	
	// 确保日志目录存在
	if err := pathutil.EnsureDir(pathutil.GetDataPath(pathutil.LOG_DIR)); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 确保日志文件存在
	if err := pathutil.EnsureFile(logPath); err != nil {
		log.Printf("创建日志文件失败: %v", err)
		return nil, nil
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var writer = file
	var flags = log.LstdFlags

	if config.TaskLogFormat && name != "main.log" {
		// 对于任务日志，使用自定义writer
		writer = &taskLogWriter{file: file}
		flags = 0 // 不需要标准日志前缀
	}

	logger := log.New(writer, "", flags)
	return logger, file
}

// CleanLogs 清理过期日志
func CleanLogs(cleanDays int) error {
	if cleanDays <= 0 {
		return fmt.Errorf("清理天数必须大于0")
	}

	logDir := pathutil.GetDataPath(pathutil.LOG_DIR)
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("读取日志目录失败: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -cleanDays)
	
	for _, file := range files {
		// 跳过main.log
		if file.Name() == "main.log" {
			continue
		}

		filePath := filepath.Join(logDir, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("获取文件信息失败[%s]: %v", file.Name(), err)
			continue
		}

		// 如果文件修改时间早于截止时间，删除文件
		if fileInfo.ModTime().Before(cutoffTime) {
			if err := os.Remove(filePath); err != nil {
				log.Printf("删除文件失败[%s]: %v", file.Name(), err)
				continue
			}
			log.Printf("已删除过期日志: %s", file.Name())
		}
	}

	return nil
}
