package static

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
)

//go:embed assets
var Assets embed.FS

// 打开文件
func OpenFile(filePath string) ([]byte, error) {
	// "assets/sql/task.sql"   文件名要以内嵌的文件夹名称开始
	entries, err := Assets.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		log.Println("打开内嵌文件失败" + filePath)
		return nil, err
	}
	return entries, nil
}

// 获取文件列表
func FileList(dirPath string) ([]map[string]string, error) {
	// "assets/sql/task.sql"   文件名要以内嵌的文件夹名称开始
	entries, err := Assets.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		log.Println("打开内嵌文件夹失败" + dirPath)
		return nil, err
	}
	var arr []map[string]string
	for _, entry := range entries {
		// fs.DirEntry的Info接口会返回fs.FileInfo，这东西被从os移动到了io/fs，接口本身没有变化
		item := make(map[string]string)
		item["name"] = entry.Name()
		item["path"] = dirPath + "/" + entry.Name()
		arr = append(arr, item)
	}
	return arr, nil
}

//将内嵌资源开放静态资源 抛出fs文件系统

func StaticFS() (fs.FS, error) {
	filesys, err := fs.Sub(Assets, "assets/www")
	if err != nil { //读取内嵌文件失败
		return nil, err
	}
	//这里可以打印调试查看内嵌文件列表
	_, err = fs.ReadDir(filesys, ".")
	if err != nil {
		log.Println("打开内嵌文件夹失败", err)
	}
	// fmt.Println(entries, 666)
	// var arr []map[string]string
	// for _, entry := range entries {
	// fs.DirEntry的Info接口会返回fs.FileInfo，这东西被从os移动到了io/fs，接口本身没有变化
	// item := make(map[string]string)
	// item["name"] = entry.Name()
	// fmt.Println(entry)
	// item["path"] = dirPath + "/" + entry.Name()
	// arr = append(arr, item)
	// }
	return filesys, nil
}
