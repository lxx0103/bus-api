package auth

import "bus-api/core/request"

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type SigninResponse struct {
	Token string `json:"token"`
	User  AdminUserResponse
}

type AdminUserResponse struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Role     string `db:"role" json:"role"`
	Status   int    `db:"status" json:"status"`
}

type AdminUserFilter struct {
	Username string `form:"username" binding:"omitempty"`
	request.PageInfo
}

type AdminUserID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type AdminPasswordUpdate struct {
	Password string `json:"password" binding:"required,min=6"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type WxSigninRequest struct {
	Code     string `json:"code" binding:"required"`
	Identity string `json:"identity" binding:"omitempty"`
}

type WechatCredential struct {
	OpenID     string `json:"openid" binding:"required"`
	SessionKey string `json:"session_key" binding:"required"`
	UnionID    string `json:"union_id"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WxSigninResponse struct {
	Token string `json:"token"`
	User  WxUserResponse
}

type WxUserResponse struct {
	ID       int64  `db:"id" json:"id"`
	OpenID   string `db:"open_id" json:"open_id"`
	Name     string `db:"name" json:"name"`
	School   string `db:"school" json:"school"`
	Grade    string `db:"grade" json:"grade"`
	Class    string `db:"class" json:"class"`
	Identity string `db:"identity" json:"identity"`
	Status   int    `db:"status" json:"status"`
}

type WxSignupRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Grade    string `json:"grade" binding:"omitempty"`
	Class    string `json:"class" binding:"omitempty"`
	Identity string `json:"identity" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type StaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type StaffFilter struct {
	Username string `form:"username" binding:"omitempty"`
	request.PageInfo
}

type StaffID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type StaffPasswordUpdate struct {
	Password string `json:"password" binding:"required,min=6"`
	UserID   int64  `json:"user_id" swaggerignore:"true"`
}

type StaffResponse struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Status   int    `db:"status" json:"status"`
}

type StaffSigninResponse struct {
	Token string `json:"token"`
	User  StaffResponse
}
