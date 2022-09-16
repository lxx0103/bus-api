package setting

import "time"

type Unit struct {
	ID             int64     `db:"id" json:"id"`
	UnitType       string    `db:"unit_type" json:"unit_type"`
	UnitID         string    `db:"unit_id" json:"unit_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
