package auth

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type authQuery struct {
	conn *sqlx.DB
}

func NewAuthQuery(connection *sqlx.DB) *authQuery {
	return &authQuery{
		conn: connection,
	}
}

func (r *authQuery) GetAdminUserByUsername(username string) (*AdminUser, error) {
	var user AdminUser
	err := r.conn.Get(&user, `
		SELECT *
		FROM u_admin_users
		WHERE username = ? AND status > 0
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authQuery) GetAdminUserCount(filter AdminUserFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Username; v != "" {
		where, args = append(where, "username like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM u_admin_users
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetAdminUserList(filter AdminUserFilter) (*[]AdminUserResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Username; v != "" {
		where, args = append(where, "username like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var users []AdminUserResponse
	err := r.conn.Select(&users, `
		SELECT id, username, role, status
		FROM u_admin_users
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *authQuery) GetAdminUserByID(id int64) (*AdminUserResponse, error) {
	var user AdminUserResponse
	err := r.conn.Get(&user, `
		SELECT id, username, role, status
		FROM u_admin_users
		WHERE id = ? AND status > 0
	`, id)
	return &user, err
}

func (r *authQuery) GetWxUserByOpenID(openID string) (*WxUserResponse, error) {
	var user WxUserResponse
	err := r.conn.Get(&user, `
		SELECT id, open_id, name, school, grade, class, identity, status
		FROM u_wx_users
		WHERE open_id = ? AND status > 0
	`, openID)
	return &user, err
}

func (r *authQuery) GetStaffCount(filter StaffFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Username; v != "" {
		where, args = append(where, "username like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM u_staffs
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetStaffList(filter StaffFilter) (*[]StaffResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Username; v != "" {
		where, args = append(where, "username like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var users []StaffResponse
	err := r.conn.Select(&users, `
		SELECT id, username, status
		FROM u_staffs
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *authQuery) GetStaffByID(id int64) (*StaffResponse, error) {
	var user StaffResponse
	err := r.conn.Get(&user, `
		SELECT id, username, status
		FROM u_staffs
		WHERE id = ? AND status > 0
	`, id)
	return &user, err
}

func (r *authQuery) GetStaffByUsername(username string) (*Staff, error) {
	var user Staff
	err := r.conn.Get(&user, `
		SELECT *
		FROM u_staffs
		WHERE username = ? AND status > 0
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
