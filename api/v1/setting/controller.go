package setting

import (
	"bus-api/core/response"
	"bus-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 单位列表
// @Id 301
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "单位名称"
// @Param unit_type query string false "单位类型（weight/length/custom)"
// @Success 200 object response.ListRes{data=[]UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units [GET]
func GetUnitList(c *gin.Context) {
	var filter UnitFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	settingService := NewSettingService()
	count, list, err := settingService.GetUnitList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建单位
// @Id 302
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param unit_info body UnitNew true "单位信息"
// @Success 200 object response.SuccessRes{data=UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units [POST]
func NewUnit(c *gin.Context) {
	var info UnitNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	settingService := NewSettingService()
	new, err := settingService.NewUnit(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新单位
// @Id 303
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Param unit_info body UnitNew true "单位信息"
// @Success 200 object response.SuccessRes{data=Unit} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [PUT]
func UpdateUnit(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info UnitNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	settingService := NewSettingService()
	new, err := settingService.UpdateUnit(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取单位
// @Id 304
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Success 200 object response.SuccessRes{data=UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [GET]
func GetUnitByID(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	unit, err := settingService.GetUnitByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, unit)

}

// @Summary 根据ID删除单位
// @Id 305
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [DELETE]
func DeleteUnit(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteUnit(uri.ID, claims.UserName)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
