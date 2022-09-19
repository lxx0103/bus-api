package qrcode

import (
	"database/sql"
	"time"
)

type qrcodeRepository struct {
	tx *sql.Tx
}

func NewQrcodeRepository(transaction *sql.Tx) *qrcodeRepository {
	return &qrcodeRepository{
		tx: transaction,
	}
}

func (r *qrcodeRepository) GetQrcodeByCode(code string) (*QrcodeResponse, error) {
	var res QrcodeResponse
	row := r.tx.QueryRow(`SELECT id, code, user_id, user_name, expire_time, status FROM q_qrcodes WHERE code = ? AND status = 1 AND expire_time > ? LIMIT 1`, code, time.Now())
	err := row.Scan(&res.ID, &res.Code, &res.UserID, &res.UserName, &res.ExpireTime, &res.Status)
	return &res, err
}

func (r *qrcodeRepository) CreateQrcode(info Qrcode) error {
	_, err := r.tx.Exec(`
		INSERT INTO q_qrcodes
		(
			user_id,
			user_name,
			code,
			expire_time,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.UserID, info.UserName, info.Code, info.ExpireTime, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *qrcodeRepository) DeleteUserQrcode(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update q_qrcodes SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE user_id = ?
		AND status = 1
	`, time.Now(), byUser, id)
	return err
}

func (r *qrcodeRepository) ScanQrcode(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update q_qrcodes SET
		status = 2,
		updated = ?,
		updated_by = ?
		WHERE id = ?
		AND status = 1
	`, time.Now(), byUser, id)
	return err
}

func (r *qrcodeRepository) CreateHistory(info ScanHistory) error {
	_, err := r.tx.Exec(`
		INSERT INTO q_scan_historys
		(
			code,
			user,
			user_id,
			by_user,
			by_user_id,
			scan_time,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Code, info.User, info.UserID, info.ByUser, info.ByUserID, info.ScanTime, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *qrcodeRepository) CheckQrCodePeriod(id int64, period time.Time) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM q_scan_historys WHERE user_id = ? AND scan_time > ?", id, period)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}
