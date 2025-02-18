package serve

import (
	"log"
	"xuanwu/config"
	r "xuanwu/gin/response"
	"xuanwu/lib"
	"xuanwu/lib/pathutil"
	"xuanwu/xuanwu"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

// UserProfile 用户配置结构体
type UserProfile struct {
	Username         string `json:"username,omitempty"`          // 用户名
	Password         string `json:"password,omitempty"`          // 密码(SHA256)
	OldPassword      string `json:"old_password,omitempty"`      // 旧密码(SHA256)
	CookieExpireDays int    `json:"cookie_expire_days,omitempty"` // Cookie过期天数
	LogCleanDays     int    `json:"log_clean_days,omitempty"`     // 日志清理天数
}

// UserInfo 用户基本信息
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 全局配置缓存
var (
	globalCookieExpireDays = 30 // 默认30天
	globalLogCleanDays     = 7  // 默认7天
)

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig() {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		return
	}
	
	if days := cfg.Get("cookie_expire_days").Int(); days > 0 {
		globalCookieExpireDays = int(days)
	}
	
	if days := cfg.Get("log_clean_days").Int(); days > 0 {
		globalLogCleanDays = int(days)
	}
}

// GetCookieExpireDays 获取当前Cookie过期天数
func GetCookieExpireDays() int {
	return globalCookieExpireDays
}

// GetLogCleanDays 获取当前日志清理天数
func GetLogCleanDays() int {
	return globalLogCleanDays
}

/* 从json文件中获取用户信息 公共方法 */
func GetUserInfo() UserInfo {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		log.Println("读取配置文件出错")
		return UserInfo{}
	}
	userInfo := UserInfo{
		Username: cfg.Get("username").String(),
		Password: cfg.Get("password").String(),
	}
	if userInfo.Username == "" { //没有查到用户数据
		return UserInfo{}
	}
	return userInfo
}

// HandlerGetUserProfile 获取用户配置
func (p *ApiData) HandlerGetUserProfile(c *gin.Context) {
	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		r.ErrMesage(c, "读取配置文件失败")
		return
	}

	profile := UserProfile{
		Username:         cfg.Get("username").String(),
		CookieExpireDays: int(cfg.Get("cookie_expire_days").Int()),
		LogCleanDays:     int(cfg.Get("log_clean_days").Int()),
	}

	if profile.Username == "" {
		r.ErrMesage(c, "获取用户信息失败")
		return
	}

	r.OkData(c, gin.H{
		"profile": profile,
		"version": config.Version,
	})
}

// HandlerUpdateUserProfile 更新用户配置
func (p *ApiData) HandlerUpdateUserProfile(c *gin.Context) {
	var req UserProfile
	if err := c.BindJSON(&req); err != nil {
		r.ErrMesage(c, "请求参数错误")
		return
	}

	// 记录请求参数（不包含密码）
	log.Printf("更新用户配置请求: username=%v, cookie_days=%v, log_days=%v", 
		req.Username, req.CookieExpireDays, req.LogCleanDays)

	cfg, err := config.ReadConfigFileToJson()
	if err != nil {
		r.ErrMesage(c, "读取配置文件失败")
		return
	}

	jsonStr := cfg.Raw
	needResetToken := false

	// 更新用户名
	if req.Username != "" {
		// 验证用户名加解密
		encryptedUsername, _ := lib.EncryptByAes([]byte(req.Username))
		decryptedUsername, err := lib.DecryptByAes(encryptedUsername)
		if err != nil || string(decryptedUsername) != req.Username {
			r.ErrMesage(c, "用户名格式错误")
			return
		}
		jsonStr, _ = sjson.Set(jsonStr, "username", req.Username)
		needResetToken = true
		log.Printf("用户名已更新为: %s", req.Username)
	}

	// 更新密码
	if req.Password != "" && req.OldPassword != "" {
		currentPass := cfg.Get("password").String()
		log.Printf("密码验证: 当前密码长度=%d, 提供的旧密码长度=%d", 
			len(currentPass), len(req.OldPassword))

		if currentPass != req.OldPassword {
			log.Printf("密码验证失败: 旧密码不匹配")
			r.ErrMesage(c, "旧密码错误")
			return
		}

		if req.Password == req.OldPassword {
			log.Printf("密码验证失败: 新密码与旧密码相同")
			r.ErrMesage(c, "新密码不能与旧密码相同")
			return
		}

		jsonStr, _ = sjson.Set(jsonStr, "password", req.Password)
		needResetToken = true
		log.Printf("密码已更新")
	} else if req.Password != "" {
		log.Printf("密码更新失败: 未提供旧密码")
		r.ErrMesage(c, "请提供旧密码")
		return
	}

	// 更新Cookie过期天数
	if req.CookieExpireDays > 0 {
		jsonStr, _ = sjson.Set(jsonStr, "cookie_expire_days", req.CookieExpireDays)
		globalCookieExpireDays = req.CookieExpireDays
		log.Printf("Cookie过期天数已更新为: %d", req.CookieExpireDays)
	}

	// 更新日志清理天数
	if req.LogCleanDays > 0 {
		jsonStr, _ = sjson.Set(jsonStr, "log_clean_days", req.LogCleanDays)
		globalLogCleanDays = req.LogCleanDays
		// 更新系统任务中的清理天数
		xuanwu.UpdateLogCleanDays(req.LogCleanDays)
		log.Printf("日志清理天数已更新为: %d", req.LogCleanDays)
	}

	// 写入配置文件
	configPath := pathutil.GetDataPath("config.json")
	if err := config.WriteConfigFile(configPath, []byte(jsonStr)); err != nil {
		log.Printf("配置文件写入失败: %v", err)
		r.ErrMesage(c, "配置文件写入失败")
		return
	}

	// 如果修改了用户名或密码，强制用户重新登录
	if needResetToken {
		p.ClearUserToken(c)
		log.Printf("用户token已清除，需要重新登录")
	}

	r.OkMesage(c, "更新成功")
} 