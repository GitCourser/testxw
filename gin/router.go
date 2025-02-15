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
	// RootRoute.Use(installHandler()) //使用中间件进行全局用户认证
	// RootRoute.Use(Cors())
	RootRoute.Use(p.CookieHandler()) //使用中间件进行全局用户认证

	routeApi := RootRoute.Group("/api") //  api接口总路由
	filesys, err := static.StaticFS()
	if err != nil {
		log.Println("加载后台文件失败,web服务停止")
		return
	}
	RootRoute.StaticFS("/admin", http.FS(filesys))

	routeAdmin := routeApi.Group("/user") // 用户数据接口
	routeAdmin.GET("/info", p.HandlerUserInfo)
	routeAdmin.GET("/list", p.AtestHandle)
	routeAdmin.POST("/update", p.HandlerUpdateUserInfo)
	routeAdmin.POST("/updatepass", p.HandlerUpdatePass)

	routeAuth := routeApi.Group("/auth") // 用户数据接口
	routeAuth.POST("/login", p.LoginHandle)
	routeAuth.GET("/logout", p.LogoutHandler)

	routeSystem := routeApi.Group("/system") // 系统配置接口
	routeSystem.GET("/admin", p.AtestHandle)
	routeSystem.GET("/info", p.AtestHandle)

	routeCron := routeApi.Group("/cron") // 定时任务接口
	/* 任务源 */
	routeCron.GET("/alllsit", cron.HandlerAllTaskList) //获取列表
	routeCron.GET("/delete", cron.HandlerDeleteTask)   //删除源任务
	routeCron.POST("/add", cron.HandlerAddTask)        //添加任务源
	routeCron.POST("/update", cron.HandlerUpdateTask)  //校验时间表达式
	/* 运行中任务 */
	routeCron.GET("/list", cron.HandlerRunTaskList)    //获取运行中的任务列表
	routeCron.GET("/enable", cron.HandlerEnableTask)   //启用任务
	routeCron.GET("/disable", cron.HandlerDisableTask) //禁用任务
	routeCron.GET("/valid", cron.Valid)                //校验时间表达式
	routeCron.POST("/execute", cron.HandlerExecuteTask) //立即执行任务
	/* 运行日志 */
	routeCron.GET("/log", cron.HandlerAllLogList)         //获取列表
	routeCron.GET("/dellog", cron.HandlerDeleteLog)       //删除日志
	routeCron.GET("/dellogall", cron.HandlerDeleteAllLog) //删除日志
	routeCron.GET("/getlog", cron.HandlerGetLog)          //获取日志
	routeCron.GET("/downlog", cron.HandlerDownloadFile)   //获取日志

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

	//动态注册插件路由
	if p.AddApi != nil {
		for key, value := range p.AddApi {
			fmt.Println(key, value, 8888)
			RootRoute.GET(value, func(c *gin.Context) {
				c.String(http.StatusOK, "Welcome Gin Server")
			})
		}
	}

	fmt.Println("This service runs on port :" + p.Port)
	RootRoute.Run(":" + p.Port)
}
