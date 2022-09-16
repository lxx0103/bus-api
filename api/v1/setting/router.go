package setting

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/units", NewUnit)
	g.GET("/units", GetUnitList)
	g.PUT("/units/:id", UpdateUnit)
	g.GET("/units/:id", GetUnitByID)
	g.DELETE("/units/:id", DeleteUnit)
}
