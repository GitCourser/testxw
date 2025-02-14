package xuanwu

import (
	"fmt"
	"io"
	"log"
	"os"
	mylog "xuanwu/log"

	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

// 定义全局定时任务
var C *cron.Cron

// 任务信息结构体
type TaskInfo struct {
	Name        string   `json:"name"`
	Times       []string `json:"times"`  // 支持多个定时时间
	WorkDir     string   `json:"workdir"` // 工作目录
	Exec        string   `json:"exec"`
	Isrun       string   `json:"isrun"` //启动执行2
	Writer      io.Writer
	Log         *log.Logger
	System      bool
	Func        func() // 系统任务函数
	Callback    string
}

// 定时id和任务的映射表
var TaskData = map[cron.EntryID]TaskInfo{}

// 定时任务
func CronInit(cfg gjson.Result) {
	tasks := cfg.Get("task")
	C = cron.New(
		cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))

	tasks.ForEach(func(key, value gjson.Result) bool { //添加用户自定义任务
		isrun := value.Get("isrun").String()
		if isrun != "2" { //启动时候是否执行
			return true
		}
		TaskData := TaskInfo{
			Name: value.Get("name").String(),
			Times: func() []string {
				var times []string
				for _, t := range value.Get("times").Array() {
					times = append(times, t.String())
				}
				return times
			}(),
			WorkDir: value.Get("workdir").String(),
			Exec: value.Get("exec").String(),
		}
		AddRunFunc(TaskData)
		return true
	})

	// 遍历系统任务切片中的每一项
	for _, item := range SystemTask {
		if item.Isrun != "2" { //启动时候是否执行
			continue
		}
		AddRunFunc(item)
	}

	C.Start()
	defer C.Stop()
	select {}
}

/* 根据任务类型,添加任务
* name 任务名称
* times 定时时间数组
* exec 执行内容
* workDir 工作目录
 */
func AddRunFunc(TaskInfo TaskInfo) {
	logname := fmt.Sprintf("%s.log", TaskInfo.Name)
	if TaskInfo.System {
		logname = ""
	}
	
	// 初始化日志
	log, writer := mylog.LogInit(logname)
	TaskInfo.Writer = writer
	TaskInfo.Log = log
	
	// 遍历时间数组,为每个时间创建定时任务
	for _, timeStr := range TaskInfo.Times {
		// 添加定时任务
		id, err := C.AddFunc(timeStr, func() {
			if err := ExecTask(TaskInfo.Exec, TaskInfo.WorkDir, log); err != nil {
				log.Printf("任务执行失败: %v\n", err)
			}
		})
		
		if err != nil {
			log.Printf("添加定时任务失败[%s]: %v\n", timeStr, err)
			continue
		}
		
		// 保存到任务映射表
		TaskData[id] = TaskInfo
	}
}

func OneRunFunc(TaskInfo TaskInfo) {
	// 初始化日志
	os.Remove("data/logs/run-task-test.log")
	log, _ := mylog.LogInit("run-task-test.log")
	
	// 执行任务
	if err := ExecTask(TaskInfo.Exec, TaskInfo.WorkDir, log); err != nil {
		log.Printf("任务执行失败: %v\n", err)
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
