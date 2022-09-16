package setting

import (
	"bus-api/core/database"
	"bus-api/core/queue"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/streadway/amqp"
)

type NewOrganizationCreated struct {
	OrganizationID string `json:"organization_id"`
	Owner          string `json:"owner"`
	OwnerEmail     string `json:"owner_email"`
	Password       string `json:"password"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("CreateOrganizationUnits", "NewOrganizationCreated", CreateOrganizationUnits)
}

func CreateOrganizationUnits(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var event NewOrganizationCreated
	err := json.Unmarshal(d.Body, &event)
	if err != nil {
		fmt.Println(err)
		return false
	}
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		// msg := "begin transaction error"
		return false
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	var unit1, unit2, unit3, unit4, unit5, unit6 Unit
	unit1.UnitID = "unit-" + xid.New().String()
	unit1.UnitType = "length"
	unit1.OrganizationID = event.OrganizationID
	unit1.Name = "CM"
	unit1.Status = 1
	unit1.Created = time.Now()
	unit1.CreatedBy = "SIGNUP"
	unit1.Updated = time.Now()
	unit1.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit1)
	if err != nil {
		fmt.Println(err)
		return false
	}
	unit2.UnitID = "unit-" + xid.New().String()
	unit2.UnitType = "length"
	unit2.OrganizationID = event.OrganizationID
	unit2.Name = "IN"
	unit2.Status = 1
	unit2.Created = time.Now()
	unit2.CreatedBy = "SIGNUP"
	unit2.Updated = time.Now()
	unit2.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit2)
	if err != nil {
		fmt.Println(err)
		return false
	}
	unit3.UnitID = "unit-" + xid.New().String()
	unit3.UnitType = "weight"
	unit3.OrganizationID = event.OrganizationID
	unit3.Name = "KG"
	unit3.Status = 1
	unit3.Created = time.Now()
	unit3.CreatedBy = "SIGNUP"
	unit3.Updated = time.Now()
	unit3.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit3)
	if err != nil {
		fmt.Println(err)
		return false
	}
	unit4.UnitID = "unit-" + xid.New().String()
	unit4.UnitType = "weight"
	unit4.OrganizationID = event.OrganizationID
	unit4.Name = "G"
	unit4.Status = 1
	unit4.Created = time.Now()
	unit4.CreatedBy = "SIGNUP"
	unit4.Updated = time.Now()
	unit4.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit4)
	if err != nil {
		fmt.Println(err)
		return false
	}
	unit5.UnitID = "unit-" + xid.New().String()
	unit5.UnitType = "weight"
	unit5.OrganizationID = event.OrganizationID
	unit5.Name = "OZ"
	unit5.Status = 1
	unit5.Created = time.Now()
	unit5.CreatedBy = "SIGNUP"
	unit5.Updated = time.Now()
	unit5.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit5)
	if err != nil {
		fmt.Println(err)
		return false
	}
	unit6.UnitID = "unit-" + xid.New().String()
	unit6.UnitType = "weight"
	unit6.OrganizationID = event.OrganizationID
	unit6.Name = "LB"
	unit6.Status = 1
	unit6.Created = time.Now()
	unit6.CreatedBy = "SIGNUP"
	unit6.Updated = time.Now()
	unit6.UpdatedBy = "SIGNUP"
	err = repo.CreateUnit(unit6)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}
