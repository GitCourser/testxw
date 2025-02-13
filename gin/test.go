package serve

import (
	r "xuanwu/gin/response"

	"github.com/gin-gonic/gin"
)

func (p *ApiData) AtestHandle(c *gin.Context) {
	// db := sqlmodel.InitModel()
	// res, _ := db.GetUserInfo("admin")
	// fmt.Println(res, "ressss")
	// // data := p.Cookie
	// static.OpenFile("assets/www/123.txt")
	// c.JSON(200, res)
	// OkMesage(c, res)
	routeTest := p.RootRoute.Group("/test")
	routeTest.GET("/test", TestHandle)
	r.OkMesage(c, "设置成功")
}

func TestHandle(c *gin.Context) {
	// db := sqlmodel.InitModel()
	// res, _ := db.GetUserInfo("admin")
	// fmt.Println(res, "ressss")
	// // data := p.Cookie
	// static.OpenFile("assets/www/123.txt")
	// c.JSON(200, res)
}
