package user

import (
	"bus-api/core/config"
	"bus-api/core/response"
	"bus-api/service"
	"errors"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 小程序用户列表
// @Id 301
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "姓名"
// @Param class query string false "班级"
// @Param grade query string false "年级"
// @Param identity query string false "身份证号"
// @Param role query string false "角色（学生，员工）"
// @Success 200 object response.ListRes{data=[]auth.WxUserResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers [GET]
func GetWxUserList(c *gin.Context) {
	var filter WxUserFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	userService := NewUserService()
	count, list, err := userService.GetWxUserList(filter)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建小程序用户
// @Id 302
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param wxUser_info body WxUserNew true "小程序用户信息"
// @Success 200 object response.SuccessRes{data=WxUserResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers [POST]
func NewWxUser(c *gin.Context) {
	var info WxUserNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	userService := NewUserService()
	err := userService.NewWxUser(info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "创建成功")
}

// @Summary 根据ID更新小程序用户
// @Id 303
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "小程序用户ID"
// @Param wxUser_info body WxUserNew true "小程序用户信息"
// @Success 200 object response.SuccessRes{data=WxUser} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers/:id [PUT]
func UpdateWxUser(c *gin.Context) {
	var uri WxUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	var info WxUserNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	userService := NewUserService()
	err := userService.UpdateWxUser(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "更新成功")
}

// @Summary 根据ID获取小程序用户
// @Id 304
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "小程序用户ID"
// @Success 200 object response.SuccessRes{data=WxUserResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers/:id [GET]
func GetWxUserByID(c *gin.Context) {
	var uri WxUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	userService := NewUserService()
	wxUser, err := userService.GetWxUserByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, wxUser)

}

// @Summary 根据ID删除小程序用户
// @Id 305
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "小程序用户ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers/:id [DELETE]
func DeleteWxUser(c *gin.Context) {
	var uri WxUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	userService := NewUserService()
	err := userService.DeleteWxUser(uri.ID, claims.UserID)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 批量导入微信用户
// @Id 306
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param file formData file true  "上传文件"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers/batch [POST]
func NewBatchWxUser(c *gin.Context) {
	uploaded, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	dest := config.ReadConfig("file.upload_path")
	extension := filepath.Ext(uploaded.Filename)
	if extension != ".csv" {
		response.ResponseError(c, "文件格式错误", errors.New("需要导入csv文件"))
		return
	}
	newName := time.Now().Format("20060102150405") + extension
	path := dest + newName
	err = c.SaveUploadedFile(uploaded, path)
	if err != nil {
		response.ResponseError(c, "保存文件错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewUserService()
	err = settingService.BatchUploadWxUser(path, claims.UserID)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 启用禁用小程序用户
// @Id 307
// @Tags 小程序用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "小程序用户ID"
// @Param wxUser_info body WxUserStatusNew true "小程序用户信息"
// @Success 200 object response.SuccessRes{data=WxUser} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wxusers/:id/status [POST]
func SetUserStatus(c *gin.Context) {
	var uri WxUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	var info WxUserStatusNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	userService := NewUserService()
	err := userService.UpdateWxUserStatus(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "更新成功")
}
