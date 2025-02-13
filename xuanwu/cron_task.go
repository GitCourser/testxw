package xuanwu

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	mylog "xuanwu/log"

	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

// 定义全局定时任务
var C *cron.Cron

// 任务信息结构体
type TaskInfo struct {
	Name     string `json:"name"`
	Time     string
	Type     string `json:"type"`
	Exec     string `json:"exec"`
	Isrun    string `json:"isrun"` //启动执行2
	Writer   io.Writer
	Log      *log.Logger
	System   bool
	Func     func() // 系统任务函数
	Callback string
}

// 定时id和任务的映射表
var TaskData = map[cron.EntryID]TaskInfo{}

// 定时任务
func CronInit(cfg gjson.Result) {
	tasks := cfg.Get("task")
	C = cron.New(
		// cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))

	tasks.ForEach(func(key, value gjson.Result) bool { //添加用户自定义任务

		isrun := value.Get("isrun").String()
		if isrun != "2" { //启动时候是否执行
			return true
		}
		TaskData := TaskInfo{
			Name: value.Get("name").String(),
			Time: value.Get("time").String(),
			Type: value.Get("type").String(),
			Exec: value.Get("exec").String(),
		}
		AddRunFunc(TaskData)
		// if _, err := cron.Parse(time); err != nil {
		// 	// 表达式错误
		// 	log := fmt.Sprintf(`%s时间表达式错误: %s`, name, time)
		// 	fmt.Println(log)
		// 	return true
		// }
		// 继续循环
		return true
	})

	// 遍历系统任务切片中的每一项
	for _, item := range SystemTask {
		// fmt.Println(item.Name, item.Type)
		if item.Isrun != "2" { //启动时候是否执行
			continue
		}
		AddRunFunc(item)
	}

	// //每天凌晨每天的0点、 4:00执行  //可以add多个定时任务
	// c.AddFunc("0 0 4 * * ?", func() {
	// 	my.taskCommand(`docker exec mysql /bin/bash -c 'mysqldump -uroot -pmypass vpndata > /var/lib/mysql/databackup/vpndata.sql && echo $?'`)
	// })

	C.Start()
	// 获取任务列表
	// GetCronList()
	defer C.Stop()
	select {}
}

/* 根据任务类型,添加任务
* type 任务类型
* name 任务名称
* name 时间
* exec 执行内容
 */
func AddRunFunc(TaskInfo TaskInfo) {
	logname := fmt.Sprintf("%s.log", TaskInfo.Name)
	if TaskInfo.System {
		logname = ""
	}
	// 根据type处理不同任务类型
	log, writer := mylog.LogInit(logname)
	TaskInfo.Writer = writer
	TaskInfo.Log = log
	switch TaskInfo.Type {
	case "1": //command
		// cmd := value.Get("exec").String()
		// 处理command任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskCommand(TaskInfo.Exec, log)
		})
		TaskData[id] = TaskInfo
	case "2": //shell
		// 处理shell脚本任务
		// //每天凌晨每天的0点、 4:00执行  //可以add多个定时任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
		})
		TaskData[id] = TaskInfo
	case "3": //shellcontent
		// 处理shell脚本任务
		// //每天凌晨每天的0点、 4:00执行  //可以add多个定时任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskShellContent(TaskInfo.Exec, log, nil)
		})
		TaskData[id] = TaskInfo
	case "4": //http
		// 处理http任务
		// //每天凌晨每天的0点、 4:00执行  //可以add多个定时任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskHttp(TaskInfo.Exec, log)
		})
		TaskData[id] = TaskInfo
	case "5": //python
		// 处理python脚本任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
		})
		TaskData[id] = TaskInfo
	case "6": //nodejs
		// 处理nodejs脚本任务
		id, _ := C.AddFunc(TaskInfo.Time, func() {
			taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
		})
		TaskData[id] = TaskInfo
	case "func": //执行gofunc
		if TaskInfo.Func != nil {
			id, _ := C.AddFunc(TaskInfo.Time, TaskInfo.Func)
			TaskData[id] = TaskInfo
		}
	}

}

func OneRunFunc(TaskInfo TaskInfo) {
	// 根据type处理不同任务类型
	os.Remove("data/logs/run-task-test.log")
	log, _ := mylog.LogInit("run-task-test.log")
	switch TaskInfo.Type {
	case "1": //command
		taskCommand(TaskInfo.Exec, log)
	case "2": //shell
		// 处理shell脚本任务
		taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
	case "3": //shellcontent
		// 处理shell脚本任务
		taskShellContent(TaskInfo.Exec, log, nil)
	case "4": //http
		// 处理http任务
		taskHttp(TaskInfo.Exec, log)
	case "5": //python
		// 处理python脚本任务
		taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
	case "6": //nodejs
		// 处理nodejs脚本任务
		taskShellFile(TaskInfo.Exec, TaskInfo.Type, log, nil)
	}
}

/* 校验时间表达式 */
func Validate(time string) bool {
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse(time); err != nil {
		// 表达式无效
		fmt.Println("错误")
		return false
	} else {
		fmt.Println("成功")
		return true
	}
}

/* 获取运行中的任务列表 */
func GetCronList() {
	entries := C.Entries()

	for _, entry := range entries {
		log.Print(entry.ID)
		log.Print(entry.Schedule)
		log.Println(entry.Next)
	}
}
