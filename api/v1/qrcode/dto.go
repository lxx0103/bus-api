package qrcode

import "bus-api/core/request"

type QrcodeResponse struct {
	ID         int64  `db:"id" json:"id"`
	Code       string `db:"code" json:"code"`
	UserID     int64  `db:"user_id" json:"user_id"`
	UserName   string `db:"user_name" json:"user_name"`
	ExpireTime string `db:"expire_time" json:"expire_time"`
	Status     int    `db:"status" json:"status"`
}

type ScanQrcodeNew struct {
	Code   string `json:"code" binding:"required,min=1,max=64"`
	UserID int64  `json:"user_id" swaggerignore:"true"`
}

type HistoryFilter struct {
	UserID       int64  `form:"user_id" binding:"omitempty"`
	ByUserID     int64  `form:"by_user_id" binding:"omitempty"`
	UserName     string `form:"user_name" binding:"omitempty"`
	ByUserName   string `form:"by_user_name" binding:"omitempty"`
	ScanDateFrom string `form:"scan_date_from" binding:"omitempty,datetime=2006-01-02"`
	ScanDateTo   string `form:"scan_date_to" binding:"omitempty,datetime=2006-01-02"`
	request.PageInfo
}

type ScanHistoryResponse struct {
	ID       int64  `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	User     string `db:"user" json:"user"`
	UserID   int64  `db:"user_id" json:"user_id"`
	ByUser   string `db:"by_user" json:"by_user"`
	ByUserID int64  `db:"by_user_id" json:"by_user_id"`
	ScanTime string `db:"scan_time" json:"scan_time"`
	Status   int    `db:"status" json:"status"`
}
