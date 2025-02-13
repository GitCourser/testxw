package xuanwu

import (
	"log"
	"xuanwu/request"
)

/*
执行http请求任务
*/
func taskHttp(url string, log *log.Logger) {
	data, err := request.Get(url, nil)
	if err != nil {
		log.Println(err, "请求执行出错了", url)
		return
	}
	log.Println(string(data))
}
