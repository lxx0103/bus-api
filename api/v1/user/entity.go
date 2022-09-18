package user

import "time"

type WxUser struct {
	ID             int64     `db:"id" json:"id"`
	WxUserType     string    `db:"wxUser_type" json:"wxUser_type"`
	WxUserID       string    `db:"wxUser_id" json:"wxUser_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
