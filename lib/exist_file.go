package lib

import (
	"fmt"
	"os"
)

/**
判断文件是否存在
*/

func PathFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// fmt.Println("文件存在")
		return true, nil
	}
	if os.IsNotExist(err) {
		fmt.Println("文件不存在")
		return false, err
	}
	return false, err
}
