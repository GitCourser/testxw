package mylog

import (
	"log"
	"os"
	"xuanwu/lib/pathutil"
)

func LogInit(name string) (*log.Logger, *os.File) {
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

	logger := log.New(file, "", log.LstdFlags)

	return logger, file
}
