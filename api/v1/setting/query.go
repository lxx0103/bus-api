package setting

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type settingQuery struct {
	conn *sqlx.DB
}

func NewSettingQuery(connection *sqlx.DB) *settingQuery {
	return &settingQuery{
		conn: connection,
	}
}

func (r *settingQuery) GetUnitByID(organizationID, id string) (*UnitResponse, error) {
	var unit UnitResponse
	err := r.conn.Get(&unit, "SELECT unit_id, organization_id, name, status FROM s_units WHERE organization_id = ? AND unit_id = ? AND status > 0", organizationID, id)
	return &unit, err
}

func (r *settingQuery) GetUnitCount(filter UnitFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.UnitType; v != "" {
		where, args = append(where, "unit_type = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_units
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetUnitList(filter UnitFilter) (*[]UnitResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.UnitType; v != "" {
		where, args = append(where, "unit_type = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var units []UnitResponse
	err := r.conn.Select(&units, `
		SELECT unit_id, organization_id, name, status
		FROM s_units
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &units, err
}
