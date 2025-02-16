package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/tidwall/gjson"
)

var Version = "1.0.0"

// 所有结构数据
type Tconfig struct{}

// 读取json配置文件转结构体
func ReadConfigFile() (*Tconfig, error) {
	// 配置文件地址修改为在当前程序目录中的data目录下
	runpath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	dir := filepath.Dir(runpath)
	dataDir := path.Join(dir, "data")
	file := path.Join(dataDir, "config.json")
	fmt.Println(file, "config")
	jsonByte, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("配置文件读取失败")
		return nil, err
	}
	var configObj Tconfig
	if err := json.Unmarshal(jsonByte, &configObj); err != nil {
		fmt.Println("反序列化配置文件错误>>", err)
		return nil, err
	}
	return &configObj, nil
}

// 将config文件读取到json字符串
func ReadConfigFileToJson() (gjson.Result, error) {
	// 配置文件地址修改为在当前程序目录中的data目录下
	runpath, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	dir := filepath.Dir(runpath)
	dataDir := path.Join(dir, "data")
	file := path.Join(dataDir, "config.json")
	jsonByte, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("配置文件读取失败")
		/* 配置文件不存在,创建json文件 */
		str := fmt.Sprintf(`{
			"name": "xuanwu",
			"username":"admin",
			"password":"8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
			"cookie_expire_days": 30,
			"task": []
		  }`)
		err := WriteConfigFile(file, []byte(str))
		if err != nil {
			log.Println("配置文件创建失败")
			return gjson.Parse(""), err
		}
		log.Println("配置文件创建成功")
		return gjson.Parse(str), nil
	}

	result := gjson.Parse(string(jsonByte))

	// 这里我们得到一个gjson实例
	// 后面可以在任意位置重用它获取值
	// lastName := result.Get("database").String()
	// fmt.Println(lastName)
	return result, nil
}

// 写入json到config文件
func WriteConfigFile(filePth string, data []byte) error {
	// 确保data目录存在
	dir := filepath.Dir(filePth)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("创建data目录失败")
			return err
		}
	}
	f, err := os.Create(filePth)
	if err != nil {
		fmt.Println("config文件创建失败")
		return err
	} else {
		_, err = f.Write(data) // 写入文件要字节类型[]byte(data)
		if err != nil {
			// 写入失败处理
			fmt.Println("config文件写入失败")
			return err
		}
	}
	defer f.Close()
	return nil
}
