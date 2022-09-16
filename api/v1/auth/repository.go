package auth

import (
	"database/sql"
)

type authRepository struct {
	tx *sql.Tx
}

func NewAuthRepository(transaction *sql.Tx) *authRepository {
	return &authRepository{
		tx: transaction,
	}
}

// func (r *authRepository) CreateUser(info User) (int64, error) {
// 	result, err := r.tx.Exec(`
// 		INSERT INTO u_users
// 		(
// 			user_id,
// 			organization_id,
// 			role_id,
// 			user_name,
// 			email,
// 			password,
// 			status,
// 			created,
// 			created_by,
// 			updated,
// 			updated_by
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	`, info.UserID, info.OrganizationID, info.RoleID, info.UserName, info.Email, info.Password, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

// // func (r *authRepository) GetUserByID(id int64) (*UserResponse, error) {
// // 	var res UserResponse
// // 	row := r.tx.QueryRow(`
// // 	SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
// // 	FROM users u
// // 	LEFT JOIN organizations o
// // 	ON u.organization_id = o.id
// // 	WHERE u.id = ?
// // 	`, id)
// // 	err := row.Scan(&res.ID, &res.Type, &res.Identifier, &res.OrganizationID, &res.PositionID, &res.RoleID, &res.Name, &res.Email, &res.Gender, &res.Phone, &res.Birthday, &res.Address, &res.Status, &res.OrganizationName)
// // 	if err != nil {
// // 		msg := "用户不存在:" + err.Error()
// // 		return nil, errors.New(msg)
// // 	}
// // 	return &res, nil
// // }

// func (r *authRepository) CheckConfict(authType int, identifier string) (bool, error) {
// 	var existed int
// 	row := r.tx.QueryRow("SELECT count(1) FROM users WHERE type = ? AND identifier = ?", authType, identifier)
// 	err := row.Scan(&existed)
// 	if err != nil {
// 		return true, err
// 	}
// 	return existed != 0, nil
// }
// func (r *authRepository) UpdateUser(id int64, info UserResponse, by string) error {
// 	_, err := r.tx.Exec(`
// 		Update users SET
// 		name = ?,
// 		email = ?,
// 		role_id = ?,
// 		position_id = ?,
// 		gender = ?,
// 		phone = ?,
// 		birthday = ?,
// 		address = ?,
// 		status = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Email, info.RoleID, info.PositionID, info.Gender, info.Phone, info.Birthday, info.Address, info.Status, time.Now(), by, id)
// 	if err != nil {
// 		msg := "更新失败:" + err.Error()
// 		return errors.New(msg)
// 	}
// 	return nil
// }
