package auth

import (
	"bus-api/core/response"
	"bus-api/service"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// @Summary 登录
// @Id 001
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param signin_info body SigninRequest true "登录类型"
// @Success 200 object response.SuccessRes{data=SigninResponse} 登录成功
// @Failure 400 object response.ErrorRes 内部错误
// @Failure 401 object response.ErrorRes 登录失败
// @Router /signin [POST]
func Signin(c *gin.Context) {
	var signinInfo SigninRequest
	err := c.ShouldBindJSON(&signinInfo)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	authService := NewAuthService()
	userInfo, err := authService.VerifyCredential(signinInfo)
	if err != nil {
		response.ResponseUnauthorized(c, "用户信息错误", err)
		return
	}
	var AdminUserResponse AdminUserResponse
	AdminUserResponse.ID = userInfo.ID
	AdminUserResponse.Username = userInfo.Username
	AdminUserResponse.Role = userInfo.Role
	AdminUserResponse.Status = userInfo.Status

	claims := service.CustomClaims{
		UserID:   userInfo.ID,
		UserName: userInfo.Username,
		Role:     userInfo.Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + 72000,
			Issuer:    "wms",
		},
	}
	jwtServices := service.JWTAuthService()
	generatedToken := jwtServices.GenerateToken(claims)
	var res SigninResponse
	res.Token = generatedToken
	res.User = AdminUserResponse
	response.Response(c, res)
}

// @Summary 创建后台用户
// @Id 002
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param signup_info body SignupRequest true "登录类型"
// @Success 200 object response.SuccessRes{data=string} 创建成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adminusers [POST]
func NewAdminUser(c *gin.Context) {
	var info SignupRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	authService := NewAuthService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	err = authService.CreateAdminUser(info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "创建成功")
}

// @Summary 后台用户列表
// @Id 003
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param username query string false "用户名称"
// @Success 200 object response.ListRes{data=[]AdminUserResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adminusers [GET]
func GetAdminUserList(c *gin.Context) {
	var filter AdminUserFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	authService := NewAuthService()
	count, list, err := authService.GetAdminUserList(filter)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID获取后台用户
// @Id 004
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Success 200 object response.SuccessRes{data=AdminUserResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adminusers/:id [GET]
func GetAdminUserByID(c *gin.Context) {
	var uri AdminUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	authService := NewAuthService()
	user, err := authService.GetAdminUserByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, user)

}

// @Summary 更新密码
// @Id 005
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body AdminPasswordUpdate true "用户信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adminusers/:id/passwords [PUT]
func UpdateAdminPassword(c *gin.Context) {
	var uri AdminUserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	var info AdminPasswordUpdate
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.UserID = claims.UserID
	authService := NewAuthService()
	err := authService.UpdateAdminPassword(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "内部错误", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 小程序登录
// @Id 006
// @Tags 小程序管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param signin_info body WxSigninRequest true "登录类型"
// @Success 200 object response.SuccessRes{data=WxSigninResponse} 登录成功
// @Failure 400 object response.ErrorRes 内部错误
// @Failure 401 object response.ErrorRes 登录失败
// @Router /wx/signin [POST]
func WxSignin(c *gin.Context) {
	var signinInfo WxSigninRequest
	err := c.ShouldBindJSON(&signinInfo)
	if err != nil {
		response.ResponseError(c, "数据格式错误", err)
		return
	}
	authService := NewAuthService()
	userInfo, err := authService.VerifyWechatSignin(signinInfo.Code)
	if err != nil {
		if err.Error() == "新用户" && userInfo.OpenID != "" && signinInfo.Identity != "" {
			err = authService.CreateWxUser(userInfo.OpenID, signinInfo.Identity)
			if err != nil {
				response.ResponseUnauthorized(c, "创建小程序用户失败", err)
				return
			}
			userInfo, err = authService.GetWxUserByOpenID(userInfo.OpenID)
			if err != nil {
				response.ResponseUnauthorized(c, "用户不存在", err)
				return
			}
		} else {
			response.ResponseUnauthorized(c, "用户不存在", err)
			return
		}
	}

	claims := service.CustomClaims{
		UserID:   userInfo.ID,
		UserName: userInfo.Name,
		Role:     userInfo.Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + 72000,
			Issuer:    "wms",
		},
	}
	jwtServices := service.JWTAuthService()
	generatedToken := jwtServices.GenerateToken(claims)
	var res WxSigninResponse
	res.Token = generatedToken
	res.User = *userInfo
	response.Response(c, res)
}
