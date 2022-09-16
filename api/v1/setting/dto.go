package setting

import "bus-api/core/request"

type UnitFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	UnitType string `form:"unit_type" binding:"required,oneof=weight length custom"`
	request.PageInfo
}

type UnitResponse struct {
	UnitID string `db:"unit_id" json:"unit_id"`
	Name   string `db:"name" json:"name"`
	Status int    `db:"status" json:"status"`
}

type UnitNew struct {
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Status int    `json:"status" binding:"required,oneof=1 2"`
	User   string `json:"user" swaggerignore:"true"`
}

type UnitID struct {
	ID string `uri:"id" binding:"required,min=1"`
}
