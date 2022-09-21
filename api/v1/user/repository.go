package user

import (
	"bus-api/api/v1/auth"
	"database/sql"
	"errors"
	"time"
)

type userRepository struct {
	tx *sql.Tx
}

func NewUserRepository(transaction *sql.Tx) *userRepository {
	return &userRepository{
		tx: transaction,
	}
}

func (r *userRepository) GetWxUserByID(id int64) (*auth.WxUserResponse, error) {
	var res auth.WxUserResponse
	row := r.tx.QueryRow(`SELECT id, open_id, name, grade, class, identity, IFNULL(expire_date,"1970-01-01") as expire_date, role, status FROM u_wx_users WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.OpenID, &res.Name, &res.Grade, &res.Class, &res.Identity, &res.ExpireDate, &res.Role, &res.Status)
	return &res, err
}

func (r *userRepository) CheckWxIdentityConfict(wxUserID int64, identity string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM u_wx_users WHERE id != ? AND identity = ? AND status > 0", wxUserID, identity)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *userRepository) CreateWxUser(info auth.WxUser) error {
	_, err := r.tx.Exec(`
		INSERT INTO u_wx_users
		(
			name,
			role,
			grade,
			class,
			identity,
			expire_date,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Role, info.Grade, info.Class, info.Identity, info.ExpireDate, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *userRepository) UpdateWxUser(id int64, info auth.WxUser) error {
	_, err := r.tx.Exec(`
		Update u_wx_users SET
		name = ?,
		role = ?,
		grade = ?,
		class = ?,
		identity = ?,
		expire_date = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Role, info.Grade, info.Class, info.Identity, info.ExpireDate, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *userRepository) DeleteWxUser(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update u_wx_users SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *userRepository) GetWxUserByIdentity(identity string) (*auth.WxUserResponse, error) {
	var res auth.WxUserResponse
	row := r.tx.QueryRow(`
	SELECT id, open_id, name, grade, class, identity, expire_date, role, status
	FROM u_wx_users
	WHERE identity = ?
	`, identity)
	err := row.Scan(&res.ID, &res.OpenID, &res.Name, &res.Grade, &res.Class, &res.Identity, &res.ExpireDate, &res.Role, &res.Status)
	return &res, err
}

func (r *userRepository) BatchCreateWxUser(wxUsers []auth.WxUser) error {
	for _, wxUser := range wxUsers {
		var wxUserExist int
		row := r.tx.QueryRow(`SELECT count(1) FROM u_wx_users WHERE identity = ? AND status > 0 LIMIT 1`, wxUser.Identity)
		err := row.Scan(&wxUserExist)
		if err != nil {
			msg := "检查身份证出错"
			return errors.New(msg)
		}
		if wxUserExist != 0 {
			_, err = r.tx.Exec(`
			UPDATE u_wx_users SET
			name = ?,
			role = ?,
			grade = ?,
			class = ?,
			expire_date = ?,
			status = ?,
			updated = ?,
			updated_by = ?
			WHERE identity = ?
			`, wxUser.Name, wxUser.Role, wxUser.Grade, wxUser.Class, wxUser.ExpireDate, wxUser.Status, wxUser.Updated, wxUser.UpdatedBy, wxUser.Identity)
			if err != nil {
				return err
			}
		} else {
			_, err = r.tx.Exec(`
			INSERT INTO u_wx_users
			(
				name,
				role,
				grade,
				class,
				identity,
				expire_date,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, wxUser.Name, wxUser.Role, wxUser.Grade, wxUser.Class, wxUser.Identity, wxUser.ExpireDate, wxUser.Status, wxUser.Created, wxUser.CreatedBy, wxUser.Updated, wxUser.UpdatedBy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *userRepository) UpdateWxUserStatus(id int64, info auth.WxUser) error {
	_, err := r.tx.Exec(`
		Update u_wx_users SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}
