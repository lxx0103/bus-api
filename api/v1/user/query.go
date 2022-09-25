package user

import (
	"bus-api/api/v1/auth"
	"strings"

	"github.com/jmoiron/sqlx"
)

type userQuery struct {
	conn *sqlx.DB
}

func NewUserQuery(connection *sqlx.DB) *userQuery {
	return &userQuery{
		conn: connection,
	}
}

func (r *userQuery) GetWxUserByID(id int64) (*auth.WxUserResponse, error) {
	var wxUser auth.WxUserResponse
	err := r.conn.Get(&wxUser, "SELECT id, open_id, name, grade, class, identity, IFNULL(expire_date, '1970-01-01') as expire_date, role, status FROM u_wx_users WHERE id = ? AND status > 0", id)
	return &wxUser, err
}

func (r *userQuery) GetWxUserCount(filter WxUserFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Grade; v != "" {
		where, args = append(where, "grade = ?"), append(args, v)
	}
	if v := filter.Class; v != "" {
		where, args = append(where, "class = ?"), append(args, v)
	}
	if v := filter.School; v != "" {
		where, args = append(where, "school = ?"), append(args, v)
	}
	if v := filter.Identity; v != "" {
		where, args = append(where, "identity like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM u_wx_users
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *userQuery) GetWxUserList(filter WxUserFilter) (*[]auth.WxUserResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Grade; v != "" {
		where, args = append(where, "grade = ?"), append(args, v)
	}
	if v := filter.Class; v != "" {
		where, args = append(where, "class = ?"), append(args, v)
	}
	if v := filter.School; v != "" {
		where, args = append(where, "school = ?"), append(args, v)
	}
	if v := filter.Identity; v != "" {
		where, args = append(where, "identity like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var wxUsers []auth.WxUserResponse
	err := r.conn.Select(&wxUsers, `
		SELECT id, open_id, name, school, grade, class, identity, status
		FROM u_wx_users
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &wxUsers, err
}
