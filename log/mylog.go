package mylog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"xuanwu/lib"
)

const logDir = "data/logs"

func LogInit(name string) (*log.Logger, *os.File) {
	if name == "" { //没有名称时候,返回空日志
		return log.New(os.Stdout, "", 0), nil
	}
	logPath := filepath.Join(logDir, name)

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	_, err = lib.PathFileExists(logPath)
	if err != nil {
		f, err := os.Create(logPath)
		if err != nil {
			log.Println(err.Error())
			return nil, nil
		} else {
			// _, err = f.Write([]byte("要写入的文本内容"))
			log.Println(name, "日志文件创建成功")
		}
		defer f.Close()
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// writer := io.MultiWriter(file)

	logger := log.New(file, "", log.LstdFlags)

	return logger, file

}
