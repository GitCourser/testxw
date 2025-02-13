package serve

import (
	"log"
	"os"
	"xuanwu/config"
	r "xuanwu/gin/response"
	"xuanwu/lib"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

/* 返回用户详情 */
func (p *ApiData) HandlerUserInfo(c *gin.Context) {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	user := cfg.Get("username").String()
	email := cfg.Get("email").String()
	if user == "" { //没有查到用户数据
		r.ErrMesage(c, "用户名错误")
		return
	}
	r.OkData(c, gin.H{
		"username": user,
		"email":    email,
		"version":  config.Version,
	})
}

type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/* 从json文件中获取用户信息 公共方法 */
func GetUserInfo() Info {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return Info{}
	}
	UserInfo := Info{
		Username: cfg.Get("username").String(),
		Password: cfg.Get("password").String(),
	}
	if UserInfo.Username == "" { //没有查到用户数据
		return Info{}
	}
	return UserInfo
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
更新用户信息
*/
func (p *ApiData) HandlerUpdateUserInfo(c *gin.Context) {
	//定义匿名结构体，字段与json字段对应
	var req struct {
		Username string `json:"userName"`
		Email    string `json:"userEmail"`
	}
	//绑定json和结构体
	if err := c.BindJSON(&req); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}
	if req.Username == "" || req.Email == "" {
		r.ErrMesage(c, "参数不能为空")
		return
	}

	/** 为防止修改用户名导致登录失败,测试一下加解密过程**/
	//加密
	str, _ := lib.EncryptByAes([]byte(req.Username))
	//解密
	username, err := lib.DecryptByAes(str)
	if err != nil {
		r.ErrMesage(c, "用户名格式错误,请重试")
		return
	}
	if string(username) != req.Username {
		r.ErrMesage(c, "用户名修改失败,请重试")
		return
	}

	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}

	jsonStr, _ := sjson.Set(cfg.Raw, "username", req.Username)
	jsonStr, _ = sjson.Set(jsonStr, "email", req.Email)
	err = os.WriteFile("data/config.json", []byte(jsonStr), 0644)
	if err != nil {
		r.ErrMesage(c, "修改失败,配置文件写入失败")
		return
	}
	r.OkMesage(c, "修改成功")
}

/*
更新用户密码
*/
func (p *ApiData) HandlerUpdatePass(c *gin.Context) {
	//定义匿名结构体，字段与json字段对应
	var req struct {
		Password_current string `json:"password_current"`
		Password         string `json:"password"`
	}
	//绑定json和结构体
	if err := c.BindJSON(&req); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}
	if req.Password_current == "" || req.Password == "" {
		r.ErrMesage(c, "参数不能为空")
		return
	}
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return
	}
	pass := cfg.Get("password").String()
	if pass != req.Password_current {
		r.ErrMesage(c, "旧密码错误")
		return
	}
	jsonStr, _ := sjson.Set(cfg.Raw, "password", req.Password)
	err = os.WriteFile("data/config.json", []byte(jsonStr), 0644)
	if err != nil {
		r.ErrMesage(c, "修改失败,配置文件写入失败")
		return
	}
	r.OkMesage(c, "修改成功")
}
