package user

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/wxusers", GetWxUserList)
	g.PUT("/wxusers/:id", UpdateWxUser)
	g.GET("/wxusers/:id", GetWxUserByID)
	g.DELETE("/wxusers/:id", DeleteWxUser)
	g.POST("/wxusers", NewWxUser)
	g.POST("/wxusers/batch", NewBatchWxUser)
	g.POST("/wxusers/:id/status", SetUserStatus)
	g.POST("/wxusers/:id/unbind", UnbindUser)
}
