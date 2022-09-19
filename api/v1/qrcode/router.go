package qrcode

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	// g.PUT("/wxqrcodes/:id", UpdateWxQrcode)
	// g.GET("/wxqrcodes/:id", GetWxQrcodeByID)
	// g.DELETE("/wxqrcodes/:id", DeleteWxQrcode)
	g.POST("/qrcodes", NewWxQrcode)
	g.POST("/scan", ScanQrcode)

	g.GET("/historys", GetHistoryList)
}
