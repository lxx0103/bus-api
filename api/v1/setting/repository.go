package setting

import (
	"database/sql"
	"time"
)

type settingRepository struct {
	tx *sql.Tx
}

func NewSettingRepository(transaction *sql.Tx) *settingRepository {
	return &settingRepository{
		tx: transaction,
	}
}

func (r *settingRepository) GetUnitByID(unitID, organizationID string) (*UnitResponse, error) {
	var res UnitResponse
	row := r.tx.QueryRow(`SELECT unit_id, organization_id, name, status FROM s_units WHERE unit_id = ? AND organization_id = ? AND status > 0 LIMIT 1`, unitID, organizationID)
	err := row.Scan(&res.UnitID, &res.OrganizationID, &res.Name, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckUnitConfict(unitID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_units WHERE organization_id = ? AND unit_id != ? AND name = ? AND status > 0 AND unit_type = ?", organizationID, unitID, name, "custom")
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateUnit(info Unit) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_units
		(
			unit_id,
			unit_type,
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.UnitID, info.UnitType, info.OrganizationID, info.Name, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateUnit(id string, info Unit) error {
	_, err := r.tx.Exec(`
		Update s_units SET
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE unit_id = ?
	`, info.Name, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteUnit(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_units SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE unit_id = ? AND unit_type = ?
	`, time.Now(), byUser, id, "custom")
	return err
}

func (r *settingRepository) GetItemUnitCount(item_id, organizationID string) (int, error) {
	var count int
	row := r.tx.QueryRow("SELECT count(1) FROM i_items WHERE organization_id = ? AND unit_id = ? AND status > 0 ", organizationID, item_id)
	err := row.Scan(&count)
	return count, err
}
