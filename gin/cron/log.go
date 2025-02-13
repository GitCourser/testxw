package cron

import (
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
	r "xuanwu/gin/response"

	"github.com/gin-gonic/gin"
)

// 定义全局变量 logDir
const logDir = "data/logs"

type FileInfo struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
	Size int64     `json:"size"`
}

/* 获取全部日志列表 */
func HandlerAllLogList(c *gin.Context) {
	files, _ := os.ReadDir(logDir)

	var fileInfos []FileInfo

	for _, f := range files {
		fileInfo, _ := f.Info()

		fi := FileInfo{
			Name: f.Name(),
			Date: fileInfo.ModTime(),
			Size: fileInfo.Size(),
		}

		fileInfos = append(fileInfos, fi)
	}
	r.OkData(c, fileInfos)
}

/* 删除单个日志 */
func HandlerDeleteLog(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		r.ErrMesage(c, "文件名为空")
		return
	}

	// 检测文件是否在 logDir 下
	cleanPath := HasFilePath(name)
	if cleanPath == "" {
		r.ErrMesage(c, "文件不存在")
		return
	}
	err := os.Remove(cleanPath)
	if err != nil {
		r.ErrMesage(c, name+" 删除失败!")
		return
	}
	r.OkMesage(c, "删除成功")
}

/* 删除全部日志 */
func HandlerDeleteAllLog(c *gin.Context) {
	err := os.RemoveAll(logDir)
	if err != nil {
		r.ErrMesage(c, "删除失败!")
		return
	}
	os.Mkdir(logDir, 0755) // 重新创建文件夹
	r.OkMesage(c, "删除成功")
}

/* 获取单个日志内容 */
func HandlerGetLog(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		r.ErrMesage(c, "文件名为空")
		return
	}
	data, err := os.ReadFile(filepath.Join(logDir, name))
	if err != nil {
		r.ErrMesage(c, name+" 读取失败!")
		return
	}
	r.OkData(c, string(data))
}

/* 日志文件下载 */
func HandlerDownloadFile(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		r.ErrMesage(c, "文件名为空")
		return
	}

	// 检测文件是否在 logDir 下
	cleanPath := HasFilePath(name)
	if cleanPath == "" {
		r.ErrMesage(c, "文件不存在")
		return
	}
	file, err := os.Open(cleanPath)
	if err != nil {
		r.ErrMesage(c, "读取失败"+cleanPath)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+name)
	// c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	chunkSize := 1024 * 1024 // 1MB
	chunks := math.Ceil(float64(fileSize) / float64(chunkSize))

	for i := 0; i < int(chunks); i++ {
		data := make([]byte, chunkSize)

		_, err = file.Read(data)
		if err != nil && err != io.EOF {
			c.AbortWithStatus(500)
			return
		}

		c.Data(200, "application/octet-stream", data[:len(data)])

		if err == io.EOF {
			break
		}
	}
}

/* sse推送流日志 */
// func sseHandler(c *gin.Context) {
// 	// 初始化客户端事件流
// 	c.Stream(func(w io.Writer) bool {
// 		fmt.Fprintf(w, "event: connect\ndata: \n\n")

// 		// 监听通道,广播日志事件
// 		ch := subscribeToLogs()
// 		for msg := range ch {
// 			event := fmt.Sprintf("event: log\ndata: %s\n\n", msg)
// 			fmt.Fprint(w, event)
// 		}
// 		return false
// 	})
// }

/* 校验文件路径是否存在于logs目录下,防止../
只允许访问logs目录下的文件,filepath参数必须为logs子路径,不能包含../向上级目录
*/

/* 校验文件路径是否存在于 logDir 下 */
func HasFilePath(name string) string {
	finalPath := filepath.Join(logDir, name)
	if !strings.HasPrefix(finalPath, logDir) {
		return ""
	}
	return filepath.Clean(finalPath)
}