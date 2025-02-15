package cron

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"xuanwu/config"
	r "xuanwu/gin/response"
	mycron "xuanwu/xuanwu"
	mylog "xuanwu/log"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/robfig/cron/v3"
)

// TaskInfo 完整的任务信息结构
type TaskInfo struct {
	ID      string   `json:"id"`      // 运行时ID
	Next    string   `json:"next"`    // 下次执行时间
	Name    string   `json:"name"`    // 任务名称
	Times   []string `json:"times"`   // 定时表达式
	WorkDir string   `json:"workdir"` // 工作目录
	Exec    string   `json:"exec"`    // 执行命令
	Enable  bool     `json:"enable"`  // 是否启用
	Status  string   `json:"status"`  // 运行状态：running/stopped
}

// HandlerTaskList 获取所有任务列表（包含运行状态）
func HandlerTaskList(c *gin.Context) {
	// 读取配置文件获取所有任务
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		r.ErrMesage(c, "读取配置文件失败")
		return
	}

	var taskList []TaskInfo
	tasks := cfg.Get("task")
	
	// 获取所有运行中任务的映射
	runningTasks := make(map[string]cron.EntryID)
	for id, task := range mycron.TaskData {
		runningTasks[task.Name] = id
	}
	
	// 遍历配置中的所有任务
	tasks.ForEach(func(key, value gjson.Result) bool {
		task := TaskInfo{
			Name:    value.Get("name").String(),
			Times:   func() []string {
				var times []string
				for _, t := range value.Get("times").Array() {
					times = append(times, t.String())
				}
				return times
			}(),
			WorkDir: value.Get("workdir").String(),
			Exec:    value.Get("exec").String(),
			Enable:  value.Get("enable").Bool(),
			Status:  "stopped", // 默认状态为停止
		}
		
		// 如果任务正在运行，添加运行时信息
		if id, exists := runningTasks[task.Name]; exists {
			for _, entry := range mycron.C.Entries() {
				if entry.ID == id {
					task.ID = strconv.Itoa(int(entry.ID))
					task.Next = entry.Next.Format("2006-01-02 15:04:05")
					task.Status = "running"
					break
				}
			}
		}
		
		taskList = append(taskList, task)
		return true
	})

	r.OkData(c, taskList)
}

// 执行任务请求参数
type executeTaskRequest struct {
	Name    string `json:"name" binding:"required"`    // 任务名称
	Exec    string `json:"exec" binding:"required"`    // 执行的命令
	WorkDir string `json:"workdir" binding:"required"` // 工作目录
}

/* 立即执行任务 */
func HandlerExecuteTask(c *gin.Context) {
	var req executeTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}

	// 初始化日志
	logname := fmt.Sprintf("%s-execute.log", req.Name)
	taskLog, file := mylog.LogInit(logname)
	if file != nil {
		defer file.Close()
	}

	// 执行任务
	go func() {
		if err := mycron.ExecTask(req.Exec, req.WorkDir, taskLog); err != nil {
			taskLog.Printf("任务执行失败: %v\n", err)
		}
	}()
	
	r.OkMesage(c, "任务执行请求已发送")
}

/* 启用任务 */
func HandlerEnableTask(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		r.ErrMesage(c, "任务名称不能为空")
		return
	}

	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}

	// 查找并更新任务状态
	found := false
	result := cfg.Get("task")
	for i, task := range result.Array() {
		if task.Get("name").String() == name {
			found = true
			// 更新配置文件
			jp := &JsonParams{data: cfg.Raw}
			jp.Set(fmt.Sprintf("task.%v.enable", i), true)
			err := os.WriteFile("data/config.json", []byte(jp.data), 0644)
			if err != nil {
				r.ErrMesage(c, "启用失败,配置文件写入失败")
				return
			}

			// 添加到cron
			TaskData := mycron.TaskInfo{
				Name: task.Get("name").String(),
				Times: func() []string {
					var times []string
					for _, t := range task.Get("times").Array() {
						times = append(times, t.String())
					}
					return times
				}(),
				WorkDir: task.Get("workdir").String(),
				Exec:    task.Get("exec").String(),
				Enable:  true,
			}
			mycron.AddRunFunc(TaskData)
			break
		}
	}

	if !found {
		r.ErrMesage(c, "任务不存在")
		return
	}

	r.OkMesage(c, "启用成功")
}

/* 禁用任务 */
func HandlerDisableTask(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		r.ErrMesage(c, "任务名称不能为空")
		return
	}

	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}

	// 查找并更新任务状态
	found := false
	result := cfg.Get("task")
	for i, task := range result.Array() {
		if task.Get("name").String() == name {
			found = true
			// 更新配置文件
			jp := &JsonParams{data: cfg.Raw}
			jp.Set(fmt.Sprintf("task.%v.enable", i), false)
			err := os.WriteFile("data/config.json", []byte(jp.data), 0644)
			if err != nil {
				r.ErrMesage(c, "禁用失败,配置文件写入失败")
				return
			}

			// 从cron中移除任务
			for _, e := range mycron.C.Entries() {
				if taskInfo, exists := mycron.TaskData[e.ID]; exists && taskInfo.Name == name {
					mycron.C.Remove(e.ID)
					// 关闭日志文件
					if file, ok := taskInfo.Writer.(*os.File); ok {
						file.Close()
					}
					// 停止写入日志
					taskInfo.Log.SetOutput(io.Discard)
					// 从映射表中删除
					delete(mycron.TaskData, e.ID)
				}
			}
			break
		}
	}

	if !found {
		r.ErrMesage(c, "任务不存在")
		return
	}

	r.OkMesage(c, "禁用成功")
}
