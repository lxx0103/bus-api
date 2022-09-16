package auth

import (
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

func (r *authQuery) GetUserByUsername(username string) (*AdminUser, error) {
	var user AdminUser
	err := r.conn.Get(&user, `
		SELECT *
		FROM u_users
		WHERE username = ? AND status > 0
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// func (r *authQuery) GetUserCount(filter UserFilter) (int, error) {
// 	where, args := []string{"status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Type; v == "wx" {
// 		where, args = append(where, "type = ?"), append(args, 2)
// 	}
// 	if v := filter.Type; v == "admin" {
// 		where, args = append(where, "type = ?"), append(args, 1)
// 	}
// 	if v := filter.OrganizationID; v != 0 {
// 		where, args = append(where, "organization_id = ?"), append(args, v)
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM users
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authQuery) GetUserList(filter UserFilter) (*[]UserResponse, error) {
// 	where, args := []string{"u.status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "u.name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Type; v == "wx" {
// 		where, args = append(where, "u.type = ?"), append(args, 2)
// 	}
// 	if v := filter.Type; v == "admin" {
// 		where, args = append(where, "u.type = ?"), append(args, 1)
// 	}
// 	if v := filter.OrganizationID; v != 0 {
// 		where, args = append(where, "u.organization_id = ?"), append(args, v)
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var users []UserResponse
// 	err := r.conn.Select(&users, `
// 		SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
// 		FROM users u
// 		LEFT JOIN organizations o
// 		ON u.organization_id = o.id
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &users, nil
// }
