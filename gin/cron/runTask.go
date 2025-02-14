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

	"github.com/gin-gonic/gin"
)

type cronEntry struct {
	ID      string   `json:"id"`
	Next    string   `json:"next"`
	Name    string   `json:"name"`
	Times   []string `json:"times"`
	WorkDir string   `json:"workdir"`
	Exec    string   `json:"exec"`
}

func HandlerRunTaskList(c *gin.Context) {
	var entries []cronEntry
	for _, e := range mycron.C.Entries() {
		customData := mycron.TaskData[e.ID]
		if customData.System { //过滤系统任务
			continue
		}
		entry := cronEntry{
			ID:      strconv.Itoa(int(e.ID)),
			Next:    e.Next.String(),
			Name:    customData.Name,
			Times:   customData.Times,
			WorkDir: customData.WorkDir,
			Exec:    customData.Exec,
		}
		entries = append(entries, entry)
	}

	r.OkData(c, entries)
}

/*
移除运行任务,接收任务id参数
*/
func HandlerRemoveTask(c *gin.Context) {
	id := c.Query("id") // 是 c.Request.URL.Query().Get("lastname") 的简写
	// 转换为 EntryID
	for _, e := range mycron.C.Entries() {
		// strconv.Itoa(int(e.ID))
		eid := fmt.Sprint(e.ID)
		if eid == id {
			mycron.C.Remove(e.ID)
			// 把日志文件释放掉
			TaskInfo := mycron.TaskData[e.ID]
			// 关闭日志
			if file, ok := TaskInfo.Writer.(*os.File); ok {
				file.Close()
			}
			// 停止写入日志
			TaskInfo.Log.SetOutput(io.Discard)
			// 删除完任务,再删除映射表
			delete(mycron.TaskData, e.ID)
			r.OkMesage(c, "删除成功")
			return
		}
	}
	r.ErrMesage(c, "删除失败,任务不存在")
}

/* 添加运行任务 */
func HandlerAddRunTask(c *gin.Context) {
	name := c.Query("name")
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	result := cfg.Get("task")
	for _, value := range result.Array() {
		if value.Get("name").String() == name {
			TaskData := mycron.TaskInfo{
				Name: value.Get("name").String(),
				Times: func() []string {
					var times []string
					for _, t := range value.Get("times").Array() {
						times = append(times, t.String())
					}
					return times
				}(),
				WorkDir: value.Get("workdir").String(),
				Exec:    value.Get("exec").String(),
			}
			// 添加运行任务
			mycron.AddRunFunc(TaskData)
			r.OkMesage(c, "运行成功")
			return
		}
	}
	r.ErrMesage(c, "运行失败,任务不存在")
}

/* 单次执行任务 */
func HandlerOneRunTask(c *gin.Context) {
	var task mycron.TaskInfo
	if err := c.ShouldBindJSON(&task); err != nil {
		// 处理错误
		r.ErrMesage(c, "参数错误")
	}
	go mycron.OneRunFunc(task)
	r.OkMesage(c, "执行请求已发送")
}
