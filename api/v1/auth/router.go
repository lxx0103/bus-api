package auth

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/signin", Signin)
	g.POST("/wx/signin", WxSignin)
	g.POST("/staff/signin", StaffSignin)
}

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/adminusers", NewAdminUser)
	g.GET("/adminusers", GetAdminUserList)
	g.GET("/adminusers/:id", GetAdminUserByID)
	g.PUT("/adminusers/:id/passwords", UpdateAdminPassword)
	g.POST("/clearalldata", ClearAllData)
	g.GET("/scanlimit", GetScanLimit)
	g.POST("/scanlimit", UpdateScanLimit)

	g.POST("/staffs", NewStaff)
	g.GET("/staffs", GetStaffList)
	g.GET("/staffs/:id", GetStaffByID)
	g.PUT("/staffs/:id/passwords", UpdateStaffPassword)
}
