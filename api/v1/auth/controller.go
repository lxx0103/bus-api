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
// @Tags 用户权限
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
		response.ResponseUnauthorized(c, "AuthError", err)
		return
	}
	var userResponse UserResponse
	userResponse.ID = userInfo.ID
	userResponse.Username = userInfo.Username
	userResponse.Role = userInfo.Role

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
	res.User = userResponse
	response.Response(c, res)
}

// // @Id 003
// // @Tags 用户权限
// // @Summary 用户注册
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param signup_info body SignupRequest true "登录类型"
// // @Success 200 object response.SuccessRes{data=int} 注册成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /signup [POST]
// func Signup(c *gin.Context) {
// 	var signupInfo SignupRequest
// 	err := c.ShouldBindJSON(&signupInfo)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	authID, err := authService.CreateAuth(signupInfo)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, authID)
// }

// // @Summary 根据ID更新用户
// // @Id 23
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users/:id [PUT]
// func UpdateUser(c *gin.Context) {
// 	var uri UserID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var user UserUpdate
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	user.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateUser(uri.ID, user, claims.UserID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 用户列表
// // @Id 32
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "用户名称"
// // @Param type query string false "用户类型wx/admin"
// // @Param organization_id query int false "用户组织"
// // @Success 200 object response.ListRes{data=[]UserResponse} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users [GET]
// func GetUserList(c *gin.Context) {
// 	var filter UserFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	organizationID := claims.OrganizationID
// 	authService := NewAuthService()
// 	count, list, err := authService.GetUserList(filter, organizationID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 根据ID获取用户
// // @Id 33
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users/:id [GET]
// func GetUserByID(c *gin.Context) {
// 	var uri UserID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	organizationID := claims.OrganizationID
// 	authService := NewAuthService()
// 	user, err := authService.GetUserByID(uri.ID, organizationID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, user)

// }

// // @Summary 用户列表
// // @Id 81
// // @Tags 小程序接口
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "用户名称"
// // @Success 200 object response.ListRes{data=[]Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /wx/users [GET]
// func WxGetUserList(c *gin.Context) {
// 	GetUserList(c)
// }

// // @Summary 根据ID更新用户
// // @Id 88
// // @Tags 小程序接口
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /wx/users/:id [PUT]
// func WxUpdateUser(c *gin.Context) {
// 	UpdateUser(c)
// }

// // @Summary 更新密码
// // @Id 102
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /password [POST]
// func UpdatePassword(c *gin.Context) {
// 	var info PasswordUpdate
// 	if err := c.ShouldBindJSON(&info); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	info.User = claims.Username
// 	info.UserID = claims.UserID
// 	authService := NewAuthService()
// 	err := authService.UpdatePassword(info)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "ok")
// }
