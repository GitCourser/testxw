package xuanwu

import (
	"log"
)

/*
系统任务
*/
var SystemTask = []TaskInfo{
	{
		Name:    "每周定时检测更新版本",
		Times:   []string{},
		WorkDir: "",
		Exec:    "",
		System:  true,
		Enable:  false,
	},
	{
		Name:    "定时清理日志或者文件",
		Times:   []string{},
		WorkDir: "",
		Exec:    "",
		System:  true,
		Enable:  false,
	},
	{
		Name:    "定时检测系统状态",
		Times:   []string{},
		WorkDir: "",
		Exec:    "",
		System:  true,
		Enable:  true,
	},
}
