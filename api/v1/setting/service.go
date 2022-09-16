package setting

import (
	"bus-api/core/database"
	"errors"
	"time"

	"github.com/rs/xid"
)

type settingService struct {
}

func NewSettingService() *settingService {
	return &settingService{}
}

func (s *settingService) GetUnitByID(organizationID, id string) (*UnitResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	unit, err := query.GetUnitByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *settingService) NewUnit(info UnitNew) (*UnitResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckUnitConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Unit name conflict"
		return nil, errors.New(msg)
	}
	var unit Unit
	unit.UnitID = "unit-" + xid.New().String()
	unit.OrganizationID = info.OrganizationID
	unit.Name = info.Name
	unit.UnitType = "custom"
	unit.Status = info.Status
	unit.Created = time.Now()
	unit.CreatedBy = info.User
	unit.Updated = time.Now()
	unit.UpdatedBy = info.User
	err = repo.CreateUnit(unit)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetUnitByID(unit.UnitID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetUnitList(filter UnitFilter) (int, *[]UnitResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetUnitCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetUnitList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateUnit(unitID string, info UnitNew) (*UnitResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckUnitConfict(unitID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "unit name conflict"
		return nil, errors.New(msg)
	}
	oldUnit, err := repo.GetUnitByID(unitID, info.OrganizationID)
	if err != nil {
		msg := "Unit not exist"
		return nil, errors.New(msg)
	}
	if oldUnit.OrganizationID != info.OrganizationID {
		msg := "Unit not exist"
		return nil, errors.New(msg)
	}
	var unit Unit
	unit.Name = info.Name
	unit.UpdatedBy = info.User
	unit.Updated = time.Now()
	unit.Status = info.Status
	err = repo.UpdateUnit(unitID, unit)
	if err != nil {
		msg := "update unit error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetUnitByID(unitID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteUnit(unitID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetUnitByID(unitID, organizationID)
	if err != nil {
		msg := "Unit not exist"
		return errors.New(msg)
	}
	itemUnitCount, err := repo.GetItemUnitCount(unitID, organizationID)
	if err != nil {
		msg := "get item unit count error"
		return errors.New(msg)
	}
	if itemUnitCount > 0 {
		msg := "Unit used by item can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteUnit(unitID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
