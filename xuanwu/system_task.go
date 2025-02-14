package xuanwu

import (
	"fmt"
	"log"
	"xuanwu/config"
	"github.com/tidwall/gjson"
)

/*
系统任务
*/
var SystemTask = []TaskInfo{
	{
		Name:    "定时清理日志或者文件",
		Times:   func() []string {
            results := value.Get("times").Array()
            strSlice := make([]string, 0, len(results))
            for _, res := range results {
                if res.Type == gjson.String {
                    strSlice = append(strSlice, res.String())
                }
            }
            return strSlice
        }(),
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
