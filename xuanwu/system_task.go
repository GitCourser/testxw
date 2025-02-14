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
		Name:    "定时清理日志或者文件",
		Times:   []string{"0 0 0 * * ?"},
		WorkDir: "",
		Exec:    "",
		System:  true,
		Isrun:   "1",
	},
}

/*
获取版本号
*/
func GetVersion() {
	data := fmt.Sprintf(`{"code":200,"message":"ok","data":{"version":"%s"}}`, config.Version)
	log.Println(data)
}
