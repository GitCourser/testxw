package serve

import (
	"net/http"
	"strings"
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
						/* token已经base64编码传输过程中不会变化格式,如果需要通过url传递,需要编码,解密遇到%会报错 */
						// ieIa9LUEblfNzTnH7UrtcA%3D%3D => ieIa9LUEblfNzTnH7UrtcA== //需要进行url解码
						// cookie, _ = url.QueryUnescape(cookie)
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
	// fmt.Println(lib.MD5(a), "lib.MD5(a)")
	//密码在设置时候就加密存储,传过来的参数是md5,直接比较
	if res.Password != req.Password {
		r.ErrMesage(c, "密码错误")
		return
	}
	//加密
	str, _ := lib.EncryptByAes([]byte(req.Username))
	//单位是秒。60*60表示60乘以60,即3600秒,也就是1小时。
	c.SetCookie("cookie", str, 60*60*24*30, "/", "", false, false)
	/* 这里返回提示即可, 用户数据会通过接口获取 */
	r.OkMesageData(c, "登录成功", gin.H{
		"token":  str,
		"maxAge": 1000 * 60 * 60 * 24 * 30,
	})
}

// 清除cookie 退出登录方法
func (p *ApiData) LogoutHandler(c *gin.Context) {
	/* 第一个参数为 cookie 名；
	第二个参数为 cookie 值；
	第三个参数为 cookie 有效时长，当 cookie 存在的时间超过设定时间时，cookie 就会失效，它就不再是我们有效的 cookie；单位是秒。
	第四个参数为 cookie 所在的目录；
	第五个为所在域，表示我们的 cookie 作用范围；
	第六个表示是否只能通过 https 访问；
	第七个表示 cookie 是否可以通过 js代码进行操作。 */
	c.SetCookie("cookie", "", -1, "/", "", false, false)
	data := "退出登录成功"
	r.OkMesage(c, data)
}
