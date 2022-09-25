package user

import "bus-api/core/request"

type WxUserFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	Grade    string `form:"grade" binding:"omitempty,max=64,min=1"`
	Class    string `form:"class" binding:"omitempty,max=64,min=1"`
	Identity string `form:"identity" binding:"omitempty,max=64,min=1"`
	School   string `form:"school" binding:"omitempty,max=64,min=1"`
	request.PageInfo
}

type WxUserResponse struct {
	WxUserID       string `db:"wxUser_id" json:"wxUser_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Name           string `db:"name" json:"name"`
	Status         int    `db:"status" json:"status"`
}

type WxUserNew struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Grade    string `json:"grade" binding:"omitempty,min=1,max=64"`
	Class    string `json:"class" binding:"omitempty,min=1,max=64"`
	School   string `json:"school" binding:"required,min=1,max=64"`
	Identity string `json:"identity" binding:"required,min=1,max=64"`
	UserID   int64  `json:"user" swaggerignore:"true"`
}

type WxUserID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type WxUserStatusNew struct {
	Status string `json:"status" binding:"required,oneof=active deactive"`
	UserID int64  `json:"user" swaggerignore:"true"`
}
