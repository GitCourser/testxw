package serve

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	r "xuanwu/gin/response"
	"xuanwu/lib"

	"github.com/gin-gonic/gin"
)

func (p *ApiData) CookieHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//这里做用户认证处理
		//判断请求是不是admin后台静态资源,不做权限认证
		if strings.HasPrefix(c.Request.RequestURI, "/admin") {
			c.Next()
		} else if strings.HasPrefix(c.Request.RequestURI, "/api") {
			cookie, err := c.Cookie("cookie")
			//cookie不存在,用户认证失败
			if c.FullPath() != "/api/auth/login" {
				if err != nil {
					//如果cookie为空,就获取Authorization
					if _, ok := c.Request.Header["Authorization"]; ok {
						// 存在
						cookie = c.Request.Header["Authorization"][0]
					} else {
						//除过login其他都要鉴权
						r.AuthMesage(c)
						c.Abort()
						return
					}
				}
				//解密
				username, err := lib.DecryptByAes(cookie)
				if err != nil {
					r.AuthMesage(c)
					c.Abort()
					return
				}
				p.Cookie = string(username)
				p.Token = cookie
			}
		} else {
			//重定向
			c.Redirect(http.StatusMovedPermanently, "/admin/")
		}
		// after request  请求前处理
		c.Next()
	}
}

// 用户登录方法
func (p *ApiData) LoginHandle(c *gin.Context) {
	//定义匿名结构体，字段与json字段对应
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//绑定json和结构体
	if err := c.BindJSON(&req); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}
	res := GetUserInfo()
	if res.Username != req.Username { //没有查到用户数据
		r.ErrMesage(c, "用户名错误")
		return
	}
	//密码在设置时候就加密存储,传过来的参数是sha256,直接比较
	if res.Password != req.Password {
		r.ErrMesage(c, "密码错误")
		return
	}
	
	// 生成token时加入时间戳确保唯一性
	tokenStr := fmt.Sprintf("%s_%d", req.Username, time.Now().Unix())
	//加密
	str, _ := lib.EncryptByAes([]byte(tokenStr))
	
	// 从配置获取cookie过期时间
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		r.ErrMesage(c, "读取配置失败")
		return
	}
	expireDays := cfg.Get("cookie_expire_days").Int()
	if expireDays <= 0 {
		expireDays = 30 // 默认30天
	}
	
	//设置cookie
	c.SetCookie("cookie", str, int(expireDays*24*60*60), "/", "", false, false)
	
	r.OkMesageData(c, "登录成功", gin.H{
		"token":  str,
		"maxAge": expireDays * 24 * 60 * 60,
	})
}

// 清除cookie 退出登录方法
func (p *ApiData) LogoutHandler(c *gin.Context) {
	c.SetCookie("cookie", "", -1, "/", "", false, false)
	data := "退出登录成功"
	r.OkMesage(c, data)
}
