package xuanwu

import (
	"fmt"
	"log"
	"xuanwu/config"
)

/*
系统任务
*/
var SystemTask = []TaskInfo{
	{
		Name:   "每周定时检测更新版本",
		Times:  [],
		WorkDir:   "",
		Exec:   "",
		System: true,
		Isrun:  "1",
	},
	{
		Name:   "定时清理日志或者文件",
		Times:  [],
		WorkDir:   "",
		Exec:   "",
		System: true,
		Isrun:  "1",
	},
	{
		Name:   "定时检测系统状态",
		Times:  [],
		WorkDir:   "",
		Exec:   "",
		System: true,
		Isrun:  "2",
	},
}

/*
获取版本号
*/
func GetVersion() {
	data := fmt.Sprintf(`{"code":200,"message":"ok","data":{"version":"%s"}}`, config.Version)
	log.Println(data)
}
