package serve

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"xuanwu/gin/cron"
	"xuanwu/static"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type ApiData struct {
	Cookie    string // 解析出来的username
	Token     string // 未解析的cookie,也就是token
	RootRoute *gin.Engine
	AddApi    map[string]string
	Port      string
}

func InitApi(cfg gjson.Result, addApi map[string]string) {
	ApiData := &ApiData{
		Cookie: "", //刷新token
		Port:   "4165",
	}
	ApiData.AddApi = addApi
	if cfg.Get("port").String() != "" {
		ApiData.Port = cfg.Get("port").String()
	}
	ApiData.Init()
}

func (p *ApiData) Init() {

	gin.SetMode(gin.ReleaseMode) // 关闭gin启动时路由打印
	RootRoute := gin.Default()
	p.RootRoute = RootRoute
	RootRoute.Use(p.CookieHandler()) //全局用户认证

	routeApi := RootRoute.Group("/api") //  api接口总路由
	filesys, err := static.StaticFS()
	if err != nil {
		log.Println("加载后台文件失败,web服务停止")
		return
	}
	RootRoute.StaticFS("/admin", http.FS(filesys))

	routeAdmin := routeApi.Group("/user") // 用户数据接口
	routeAdmin.GET("/profile", p.HandlerGetUserProfile)    // 获取用户配置
	routeAdmin.POST("/profile", p.HandlerUpdateUserProfile) // 更新用户配置

	routeAuth := routeApi.Group("/auth") // 用户数据接口
	routeAuth.POST("/login", p.LoginHandle)
	routeAuth.GET("/logout", p.LogoutHandler)

	routeCron := routeApi.Group("/cron") // 定时任务接口
	/* 任务源 */
	routeCron.GET("/list", cron.HandlerTaskList)    //获取任务列表（包含运行状态）
	routeCron.GET("/delete", cron.HandlerDeleteTask)   //删除源任务
	routeCron.POST("/add", cron.HandlerAddTask)        //添加任务源
	routeCron.POST("/update", cron.HandlerAddTask)     //更新任务（复用添加接口）
	/* 任务控制 */
	routeCron.GET("/enable", cron.HandlerEnableTask)   //启用任务
	routeCron.GET("/disable", cron.HandlerDisableTask) //禁用任务
	routeCron.POST("/execute", cron.HandlerExecuteTask) //立即执行任务

	// 文件管理接口
	routeFile := routeApi.Group("/file")
	routeFile.GET("/list", HandlerFileList)       // 获取文件列表
	routeFile.POST("/upload", HandlerFileUpload)  // 上传文件
	routeFile.POST("/batch-upload", HandlerBatchUpload) // 批量上传文件
	routeFile.POST("/mkdir", HandlerMkdir)       // 创建文件夹
	routeFile.GET("/download", HandlerFileDownload) // 下载文件
	routeFile.GET("/content", HandlerFileContent) // 获取文件内容
	routeFile.POST("/edit", HandlerFileEdit)     // 编辑文件
	routeFile.GET("/delete", HandlerFileDelete)  // 删除文件

	// 关键点【解决页面刷新404的问题】
	RootRoute.NoRoute(func(c *gin.Context) {
		//设置响应状态
		c.Writer.WriteHeader(http.StatusOK)
		//载入首页
		indexHTML, _ := fs.ReadFile(filesys, "index.html")
		c.Writer.Write(indexHTML)
		//响应HTML类型
		c.Writer.Header().Add("Accept", "text/html")
		//显示刷新
		c.Writer.Flush()
	})

	fmt.Println("Web 端口：" + p.Port)
	RootRoute.Run(":" + p.Port)
}
