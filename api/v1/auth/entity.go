package auth

import "time"

type AdminUser struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	Role      string    `db:"role" json:"role"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type WxUser struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	OpenID     string    `db:"open_id" json:"open_id"`
	Grade      string    `db:"grade" json:"grade"`
	Class      string    `db:"class" json:"class"`
	Identity   string    `db:"identity" json:"identity"`
	ExpireDate string    `db:"expire_date" json:"expire_date"`
	Role       string    `db:"role" json:"role"`
	Status     int       `db:"status" json:"status"`
	Created    time.Time `db:"created" json:"created"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	Updated    time.Time `db:"updated" json:"updated"`
	UpdatedBy  string    `db:"updated_by" json:"updated_by"`
}
