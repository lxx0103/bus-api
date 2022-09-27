package auth

import (
	"database/sql"
	"time"
)

type authRepository struct {
	tx *sql.Tx
}

func NewAuthRepository(transaction *sql.Tx) *authRepository {
	return &authRepository{
		tx: transaction,
	}
}

func (r *authRepository) CreateAdminUser(info AdminUser) error {
	_, err := r.tx.Exec(`
		INSERT INTO u_admin_users
		(
			username,
			role,
			password,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Username, info.Role, info.Password, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *authRepository) GetAdminUserByID(id int64) (*AdminUserResponse, error) {
	var res AdminUserResponse
	row := r.tx.QueryRow(`
	SELECT id, username, role, status
	FROM u_admin_users
	WHERE id = ?
	`, id)
	err := row.Scan(&res.ID, &res.Username, &res.Role, &res.Status)
	return &res, err
}

func (r *authRepository) CheckAdminUserConfict(username string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM u_admin_users WHERE username = ?", username)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *authRepository) UpdateAdminPassword(id int64, password, byUser string) error {
	_, err := r.tx.Exec(`
		Update u_admin_users SET
		password = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, password, time.Now(), byUser, id)
	return err
}

func (r *authRepository) CheckWxUserConfict(openID string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM u_wx_users WHERE open_id = ?", openID)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *authRepository) UpdateWxUser(id int64, info WxUser) error {
	_, err := r.tx.Exec(`
		UPDATE u_wx_users
		set open_id = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.OpenID, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *authRepository) GetWxUserByIdentity(identity string) (*WxUserResponse, error) {
	var res WxUserResponse
	row := r.tx.QueryRow(`
	SELECT id, open_id, name, school, grade, class, identity, status
	FROM u_wx_users
	WHERE identity = ? AND status > 0
	`, identity)
	err := row.Scan(&res.ID, &res.OpenID, &res.Name, &res.School, &res.Grade, &res.Class, &res.Identity, &res.Status)
	return &res, err
}

func (r *authRepository) CheckStaffConfict(username string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM u_staffs WHERE username = ?", username)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *authRepository) CreateStaff(info Staff) error {
	_, err := r.tx.Exec(`
		INSERT INTO u_staffs
		(
			username,
			password,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.Username, info.Password, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *authRepository) UpdateStaffPassword(id int64, password, byUser string) error {
	_, err := r.tx.Exec(`
		Update u_staffs SET
		password = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, password, time.Now(), byUser, id)
	return err
}

func (r *authRepository) GetStaffByID(id int64) (*StaffResponse, error) {
	var res StaffResponse
	row := r.tx.QueryRow(`
	SELECT id, username,  status
	FROM u_staffs
	WHERE id = ?
	`, id)
	err := row.Scan(&res.ID, &res.Username, &res.Status)
	return &res, err
}

func (r *authRepository) ClearAllData(byUser string) error {
	_, err := r.tx.Exec(`
		Update u_staffs SET
		status = -2,
		updated = ?,
		updated_by = ?
		WHERE status > 0
	`, time.Now(), byUser)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update u_wx_users SET
		status = -2,
		updated = ?,
		updated_by = ?
		WHERE status > 0
	`, time.Now(), byUser)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update q_scan_historys SET
		status = -2,
		updated = ?,
		updated_by = ?
		WHERE status > 0
	`, time.Now(), byUser)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		Update q_qrcodes SET
		status = -2,
		updated = ?,
		updated_by = ?
		WHERE status > 0
	`, time.Now(), byUser)
	return err
}

func (r *authRepository) UpdateScanLimit(byUser string, limit int) error {
	_, err := r.tx.Exec(`
		Update s_scan_limits SET
		limits = ?,
		updated = ?,
		updated_by = ?
		WHERE id = 1
	`, limit, time.Now(), byUser)
	return err
}
