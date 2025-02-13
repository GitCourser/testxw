package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"xuanwu/config"
	r "xuanwu/gin/response"
	mycron "xuanwu/xuanwu"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func HandlerAllTaskList(c *gin.Context) {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}

	var tempArray []interface{}
	json.Unmarshal([]byte(cfg.Get("task").Raw), &tempArray)
	r.OkData(c, tempArray)
}

// 放处理过的json字符串
type JsonParams struct {
	data string
}

/* 给json设置值 */
func (j *JsonParams) Set(key string, value interface{}) {
	j.data, _ = sjson.Set(j.data, key, value)
}

/*
添加任务
*/
func HandlerAddTask(c *gin.Context) {
	// nick := c.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值
	// 声明map结构
	var jsonData map[string]interface{}
	// 解析请求体到map
	if err := c.BindJSON(&jsonData); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}
	// 获取name参数
	name := jsonData["name"].(string)
	if name == "" {
		r.ErrMesage(c, "任务名称不能为空")
		return
	}

	taskType := jsonData["type"]
	if taskType == "" {
		r.ErrMesage(c, "任务类型不能为空")
		return
	}
	timestr := jsonData["time"]
	if ok := mycron.Validate(timestr.(string)); !ok {
		r.ErrMesage(c, "时间表达式错误")
		return
	}
	jp := &JsonParams{data: ""}
	jp.Set("name", name)
	jp.Set("type", taskType)
	jp.Set("time", timestr)
	jp.Set("exec", jsonData["exec"])
	jp.Set("isrun", jsonData["isrun"])
	jp.Set("desc", jsonData["desc"])
	jp.Set("createtime", time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}

	result := gjson.Get(cfg.Raw, "task.#.name")
	for _, isname := range result.Array() {
		if isname.String() == name {
			r.ErrMesage(c, "任务已存在")
			return
		}
		println(isname.String())
	}

	var newObj map[string]interface{}
	json.Unmarshal([]byte(jp.data), &newObj)
	value, _ := sjson.Set(cfg.Raw, "task.-1", newObj)
	err = os.WriteFile("data/config.json", []byte(value), 0644)
	if err != nil {
		r.ErrMesage(c, "添加失败,配置文件写入失败")
		return
	}
	// 这里可以返回任务列表数据
	// result := gjson.Parse(value)
	// var tempArray []interface{}
	// json.Unmarshal([]byte(result.Get("task").Raw), &tempArray)
	// r.OkData(c, tempArray)
	r.OkMesage(c, "添加成功")
}

/* 删除任务源 */
func HandlerDeleteTask(c *gin.Context) {
	name := c.Query("name") // 是 c.Request.URL.Query().Get("lastname") 的简写
	if name == "" {
		r.ErrMesage(c, "任务名称不能为空")
		return
	}
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	result := gjson.Get(cfg.Raw, "task.#.name")
	for i, isname := range result.Array() {
		if isname.String() == name {
			value, _ := sjson.Delete(cfg.Raw, fmt.Sprintf("task.%v", i))
			println(value)
			println(fmt.Sprintf("task.%v", i))
			err := os.WriteFile("data/config.json", []byte(value), 0644)
			if err != nil {
				r.ErrMesage(c, "删除失败,配置文件写入失败")
				return
			}
			r.OkMesage(c, "删除成功")
			return
		}
		// println(isname.String())
		// println(i)
	}
	r.ErrMesage(c, "删除失败,任务不存在")

}

/* 修改任务源 */
func HandlerUpdateTask(c *gin.Context) {
	// 声明map结构
	var jsonData map[string]interface{}
	// 解析请求体到map
	if err := c.BindJSON(&jsonData); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}
	name := jsonData["name"] // 是 c.Request.URL.Query().Get("lastname") 的简写
	time := jsonData["time"]
	isrun := jsonData["isrun"]
	if name == "" || time == "" || isrun == "" {
		r.ErrMesage(c, "参数不能为空")
		return
	}
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	jp := &JsonParams{data: cfg.Raw}
	result := gjson.Get(cfg.Raw, "task.#.name")
	for i, isname := range result.Array() {
		if isname.String() == name {
			jp.Set(fmt.Sprintf("task.%v.time", i), time)
			jp.Set(fmt.Sprintf("task.%v.isrun", i), isrun)
			err := os.WriteFile("data/config.json", []byte(jp.data), 0644)
			if err != nil {
				r.ErrMesage(c, "修改失败,配置文件写入失败")
				return
			}
			r.OkMesage(c, "修改成功")
			return
		}
		// println(isname.String())
		// println(i)
	}
	r.ErrMesage(c, "删除失败,任务不存在")

}

func Valid(c *gin.Context) {
	timestr := c.Query("time")
	if ok := mycron.Validate(timestr); !ok {
		r.ErrMesage(c, "时间表达式错误")
		return
	}
	r.OkMesage(c, "表达式格式正确")
}
