package qrcode

import "time"

type Qrcode struct {
	ID         int64     `db:"id" json:"id"`
	Code       string    `db:"code" json:"code"`
	UserID     int64     `db:"user_id" json:"user_id"`
	UserName   string    `db:"user_name" json:"user_name"`
	ExpireTime string    `db:"expire_time" json:"expire_time"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}

type ScanHistory struct {
	ID        int64     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	User      string    `db:"user" json:"user"`
	UserID    int64     `db:"user_id" json:"user_id"`
	ByUser    string    `db:"by_user" json:"by_user"`
	ByUserID  int64     `db:"by_user_id" json:"by_user_id"`
	ScanTime  string    `db:"scan_time" json:"scan_time"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
