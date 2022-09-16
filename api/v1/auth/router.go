package auth

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/signin", Signin)
	// g.POST("/wx/signin", WxSignin)
}

func AuthRouter(g *gin.RouterGroup) {
	// g.POST("/users", NewAdminUser)
}
