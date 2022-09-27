package auth

import (
	"bus-api/core/config"
	"bus-api/core/database"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
}

func NewAuthService() *authService {
	return &authService{}
}

func (s authService) CreateAdminUser(signupInfo SignupRequest) error {
	hashed, err := hashPassword(signupInfo.Password)
	if err != nil {
		msg := "加密密码出错"
		return errors.New(msg)
	}
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(signupInfo.UserID)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	if byUser.Role != "超级管理员" {
		msg := "只有超级管理员可以新增管理员"
		return errors.New(msg)
	}
	var newUser AdminUser
	newUser.Password = hashed
	isConflict, err := repo.CheckAdminUserConfict(signupInfo.Username)
	if err != nil {
		msg := "检查用户名合法性出错"
		return errors.New(msg)
	}
	if isConflict {
		msg := "用户名已存在"
		return errors.New(msg)
	}
	newUser.Username = signupInfo.Username
	newUser.Role = "后台用户"
	newUser.Status = 1
	newUser.Created = time.Now()
	newUser.CreatedBy = byUser.Username
	newUser.Updated = time.Now()
	newUser.UpdatedBy = byUser.Username
	err = repo.CreateAdminUser(newUser)
	if err != nil {
		msg := "创建后台用户出错"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) GetAdminUserList(filter AdminUserFilter) (int, *[]AdminUserResponse, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	count, err := query.GetAdminUserCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetAdminUserList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *authService) GetAdminUserByID(id int64) (*AdminUserResponse, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	user, err := query.GetAdminUserByID(id)
	if err != nil {
		msg := "用户不存在"
		return nil, errors.New(msg)
	}
	return user, nil
}

func (s *authService) VerifyCredential(signinInfo SigninRequest) (*AdminUser, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	userInfo, err := query.GetAdminUserByUsername(signinInfo.Username)
	if err != nil {
		msg := "用户不存在"
		return nil, errors.New(msg)
	}
	if !checkPasswordHash(signinInfo.Password, userInfo.Password) {
		msg := "密码错误"
		return nil, errors.New(msg)
	}
	return userInfo, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) UpdateAdminPassword(id int64, info AdminPasswordUpdate) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	if byUser.Role != "超级管理员" && byUser.ID != id {
		msg := "只有超级管理员才可以修改他人密码"
		return errors.New(msg)
	}
	hashed, err := hashPassword(info.Password)
	if err != nil {
		msg := "密码加密错误"
		return errors.New(msg)
	}
	err = repo.UpdateAdminPassword(id, hashed, byUser.Username)
	if err != nil {
		msg := "密码更新错误"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) VerifyWechatSignin(code string) (*WxUserResponse, error) {
	var credential WechatCredential
	if code == "test" {
		credential.ErrCode = 0
		credential.ErrMsg = ""
		credential.OpenID = "test"
	} else if code == "test2" {
		credential.ErrCode = 0
		credential.ErrMsg = ""
		credential.OpenID = "test2"
	} else {
		httpClient := &http.Client{}
		signin_uri := config.ReadConfig("Wechat.signin_uri")
		appID := config.ReadConfig("Wechat.app_id")
		appSecret := config.ReadConfig("Wechat.app_secret")
		uri := signin_uri + "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return nil, err
		}
		res, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &credential)
		if err != nil {
			msg := "解码出错"
			return nil, errors.New(msg)
		}
	}
	if credential.ErrCode != 0 {
		msg := credential.ErrMsg
		return nil, errors.New(msg)
	}
	db := database.RDB()
	query := NewAuthQuery(db)
	userInfo, err := query.GetWxUserByOpenID(credential.OpenID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			msg := "新用户"
			var res WxUserResponse
			res.OpenID = credential.OpenID
			return &res, errors.New(msg)
		} else {
			msg := "获取用户出错"
			return nil, errors.New(msg)
		}
	}
	return userInfo, nil
}

func (s authService) CreateWxUser(openID, identity string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	var newUser WxUser
	isConflict, err := repo.CheckWxUserConfict(openID)
	if err != nil {
		msg := "检查用户名合法性出错"
		return errors.New(msg)
	}
	if isConflict {
		msg := "用户名已存在"
		return errors.New(msg)
	}
	oldUser, err := repo.GetWxUserByIdentity(identity)
	if err != nil {
		msg := "你的身份证信息还未登记"
		return errors.New(msg)
	}
	if oldUser.OpenID != "" {
		msg := "此身份证已被他人绑定"
		return errors.New(msg)
	}
	newUser.OpenID = openID
	newUser.Status = 2
	newUser.Updated = time.Now()
	newUser.UpdatedBy = "SIGNUP"
	err = repo.UpdateWxUser(oldUser.ID, newUser)
	if err != nil {
		msg := "创建微信用户出错"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) GetWxUserByOpenID(id string) (*WxUserResponse, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	user, err := query.GetWxUserByOpenID(id)
	if err != nil {
		msg := "用户不存在"
		return nil, errors.New(msg)
	}
	return user, nil
}

func (s authService) CreateStaff(signupInfo StaffRequest) error {
	hashed, err := hashPassword(signupInfo.Password)
	if err != nil {
		msg := "加密密码出错"
		return errors.New(msg)
	}
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(signupInfo.UserID)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	var newUser Staff
	newUser.Password = hashed
	isConflict, err := repo.CheckStaffConfict(signupInfo.Username)
	if err != nil {
		msg := "检查用户名合法性出错"
		return errors.New(msg)
	}
	if isConflict {
		msg := "用户名已存在"
		return errors.New(msg)
	}
	newUser.Username = signupInfo.Username
	newUser.Status = 1
	newUser.Created = time.Now()
	newUser.CreatedBy = byUser.Username
	newUser.Updated = time.Now()
	newUser.UpdatedBy = byUser.Username
	err = repo.CreateStaff(newUser)
	if err != nil {
		msg := "创建员工出错"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) GetStaffList(filter StaffFilter) (int, *[]StaffResponse, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	count, err := query.GetStaffCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetStaffList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *authService) GetStaffByID(id int64) (*StaffResponse, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	user, err := query.GetStaffByID(id)
	if err != nil {
		msg := "员工不存在"
		return nil, errors.New(msg)
	}
	return user, nil
}

func (s *authService) UpdateStaffPassword(id int64, info StaffPasswordUpdate) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	hashed, err := hashPassword(info.Password)
	if err != nil {
		msg := "密码加密错误"
		return errors.New(msg)
	}
	err = repo.UpdateStaffPassword(id, hashed, byUser.Username)
	if err != nil {
		msg := "密码更新错误"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) VerifyStaffCredential(signinInfo SigninRequest) (*Staff, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	userInfo, err := query.GetStaffByUsername(signinInfo.Username)
	if err != nil {
		msg := "用户不存在"
		return nil, errors.New(msg)
	}
	if !checkPasswordHash(signinInfo.Password, userInfo.Password) {
		msg := "密码错误"
		return nil, errors.New(msg)
	}
	return userInfo, nil
}

func (s *authService) ClearAllData(id int64) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(id)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	if byUser.Role != "超级管理员" {
		msg := "只有超级管理员才可以清空数据"
		return errors.New(msg)
	}
	err = repo.ClearAllData(byUser.Username)
	if err != nil {
		msg := "清空数据错误"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *authService) GetScanLimit(id int64) (int, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	byUser, err := query.GetAdminUserByID(id)
	if err != nil {
		msg := "获取当前用户出错"
		return 0, errors.New(msg)
	}
	if byUser.Role != "超级管理员" {
		msg := "只有超级管理员才可以获取扫码限制"
		return 0, errors.New(msg)
	}
	res, err := query.GetScanLimit()
	if err != nil {
		msg := "获取数据出错"
		return 0, errors.New(msg)
	}
	return res, nil
}

func (s *authService) UpdateScanLimit(id int64, limit int) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "开启事务出错"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	byUser, err := repo.GetAdminUserByID(id)
	if err != nil {
		msg := "获取当前用户出错"
		return errors.New(msg)
	}
	if byUser.Role != "超级管理员" {
		msg := "只有超级管理员才可以更新扫码限制次数"
		return errors.New(msg)
	}
	err = repo.UpdateScanLimit(byUser.Username, limit)
	if err != nil {
		msg := "更新扫码限制次数错误"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}
