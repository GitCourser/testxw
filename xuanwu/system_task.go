package xuanwu

import (
	"log"
	"xuanwu/config"
	xwlog "xuanwu/log"
)

var systemLogger *log.Logger

func init() {
	// 初始化系统日志记录器
	logger, _ := xwlog.LogInit("main.log")
	systemLogger = logger
}

/*
系统任务
*/
var SystemTask = []TaskInfo{
	{
		Name:    "定时清理日志",
		Times:   []string{"@daily"}, // 每天0点执行
		WorkDir: "",
		Exec:    "",
		System:  true,
		Enable:  true,
		Func:    cleanLogsTask,
	},
	{
		Name:    "系统测试任务",
		Times:   []string{"@every 30s"},
		WorkDir: "",
		Exec:    "",
		System:  true,
		Enable:  true,
		Func:    systemTestTask,
	},
}

// cleanLogsTask 清理过期日志任务
func cleanLogsTask() {
	// 记录任务开始
	systemLogger.Printf("定时清理日志")
	
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return
	}

	cleanDays := cfg.Get("log_clean_days").Int()
	if cleanDays <= 0 {
		cleanDays = 7 // 默认7天
	}

	if err := xwlog.CleanLogs(int(cleanDays)); err != nil {
		log.Printf("清理日志失败: %v", err)
		return
	}
}

// 系统测试任务
func systemTestTask() {
	systemLogger.Printf("系统测试任务")
}
