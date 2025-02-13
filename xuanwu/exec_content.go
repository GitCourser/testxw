package xuanwu

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

/*

****  可以异步执行,实时结果的 ****

 */

func taskShellContent(content string, log *log.Logger, callback func()) {
	// 执行shell脚本
	res, err := execShellContent(content, log)

	if err != nil {
		log.Printf("执行脚本%s失败: %v", content, err)
	} else {
		log.Printf("执行脚本%s成功, %s", content, res)

		// 执行回调函数
		if callback != nil {
			callback()
		}
	}
}

/** 执行shell内容 **/
func execShellContent(scriptContent string, log *log.Logger) (string, error) {

	// 构建一个脚本
	var script bytes.Buffer
	script.WriteString("#!/bin/bash\n")
	scriptContent = strings.TrimSpace(scriptContent)                 // 去掉头尾空格
	scriptContent = strings.ReplaceAll(scriptContent, "\u00a0", " ") // 替换非法空格
	script.WriteString(scriptContent + "\n")
	// 创建一个bash命令
	cmd := exec.Command("bash", "-s")

	// 将脚本内容传递给标准输入
	cmd.Stdin = &script

	stdout, _ := cmd.StdoutPipe()

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Printf("Read error: %v\n", err)
					break
				}
			}

			log.Println(string(line))
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}

	return "", nil
}
