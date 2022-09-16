package service

import (
	"bus-api/core/database"
)

type rbacService struct {
}

func NewRbacService() *rbacService {
	return &rbacService{}
}

func (service *rbacService) CheckPrivilege(role_id string, route string, method string) bool {
	var res int64
	db := database.RDB()
	err := db.Get(&res, `
		SELECT count(1) FROM 
		user_role_menus rm
		LEFT JOIN user_menu_apis ma
		ON rm.menu_id = ma.menu_id
		LEFT JOIN user_apis a 
		ON ma.api_id = a.id
		WHERE rm.role_id = ?
		AND a.route = ?
		AND a.method = ?
		AND ma.enabled = 1
		AND rm.enabled = 1 	
	`, role_id, route, method)
	if err != nil {
		return false
	}
	if res < 1 {
		return false
	}
	return true
}
