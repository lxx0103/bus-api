package qrcode

import (
	"bus-api/core/response"
	"bus-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 新建二维码
// @Id 201
// @Tags 二维码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /qrcodes [POST]
func NewWxQrcode(c *gin.Context) {
	claims := c.MustGet("claims").(*service.CustomClaims)
	qrcodeService := NewQrcodeService()
	res, err := qrcodeService.NewWxQrcode(claims.UserID)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, res)
}

// @Summary 扫描二维码
// @Id 202
// @Tags 二维码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param wxQrcode_info body ScanQrcodeNew true "二维码信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /scan [POST]
func ScanQrcode(c *gin.Context) {
	var info ScanQrcodeNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	qrcodeService := NewQrcodeService()
	err := qrcodeService.ScanQrcode(info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "扫码成功")
}

// @Summary 扫码历史列表
// @Id 203
// @Tags 二维码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param user_id query string false "学生ID"
// @Param by_user_id query string false "员工ID"
// @Param user_name query string false "学生姓名"
// @Param by_user_name query string false "员工姓名"
// @Param scan_date_from query string false "扫码时间开始"
// @Param scan_date_to query string false "扫码时间结束"
// @Success 200 object response.ListRes{data=[]ScanHistoryResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /historys [GET]
func GetHistoryList(c *gin.Context) {
	var filter HistoryFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	qrcodeService := NewQrcodeService()
	count, list, err := qrcodeService.GetHistoryList(filter)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}
