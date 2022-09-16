package auth

import (
	"bus-api/core/database"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
}

func NewAuthService() *authService {
	return &authService{}
}

// func (s authService) CreateAuth(signupInfo SignupRequest) (int64, error) {
// 	hashed, err := hashPassword(signupInfo.Credential)
// 	if err != nil {
// 		msg := "hash password error"
// 		return 0, errors.New(msg)
// 	}
// 	db := database.WDB()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		msg := "begin transaction error"
// 		return 0, errors.New(msg)
// 	}
// 	defer tx.Rollback()
// 	repo := NewAuthRepository(tx)
// 	var newUser User
// 	newUser.Credential = hashed
// 	isConflict, err := repo.CheckConfict(1, signupInfo.Identifier)
// 	if err != nil {
// 		msg := "check conflict error"
// 		return 0, errors.New(msg)
// 	}
// 	if isConflict {
// 		msg := "Username conflict"
// 		return 0, errors.New(msg)
// 	}
// 	newUser.Identifier = signupInfo.Identifier
// 	newUser.AuthType = 1
// 	newUser.OrganizationID = signupInfo.OrganizationID
// 	authID, err := repo.CreateUser(newUser)
// 	if err != nil {
// 		msg := "create user error"
// 		return 0, errors.New(msg)
// 	}
// 	tx.Commit()
// 	return authID, nil
// }

// func (s *authService) VerifyWechatSignin(code string) (*WechatCredential, error) {
// 	var credential WechatCredential
// 	httpClient := &http.Client{}
// 	signin_uri := config.ReadConfig("Wechat.signin_uri")
// 	appID := config.ReadConfig("Wechat.app_id")
// 	appSecret := config.ReadConfig("Wechat.app_secret")
// 	uri := signin_uri + "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
// 	req, err := http.NewRequest("GET", uri, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, err := httpClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(body, &credential)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &credential, nil
// }

// func (s *authService) GetUserInfo(openID string, authType int, organizationID int64) (*UserResponse, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	user, err := query.GetUserByOpenID(openID)
// 	if err != nil {
// 		if err.Error() != "sql: no rows in result set" {
// 			return nil, err
// 		}
// 		if organizationID == 0 {
// 			msg := "组织ID不存在"
// 			return nil, errors.New(msg)
// 		}
// 		tx, err := db.Begin()
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer tx.Rollback()
// 		var newUser User
// 		newUser.Type = authType
// 		newUser.Identifier = openID
// 		newUser.OrganizationID = organizationID
// 		repo := NewAuthRepository(tx)
// 		userID, err := repo.CreateUser(newUser)
// 		if err != nil {
// 			return nil, err
// 		}
// 		user, err = repo.GetUserByID(userID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tx.Commit()
// 	}
// 	return user, nil
// }

func (s *authService) VerifyCredential(signinInfo SigninRequest) (*AdminUser, error) {
	db := database.RDB()
	query := NewAuthQuery(db)
	userInfo, err := query.GetUserByUsername(signinInfo.Username)
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

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func (s *authService) UpdateUser(userID int64, info UserUpdate, byUserID int64) (*UserResponse, error) {
// 	db := database.InitMySQL()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		msg := "事务开启错误" + err.Error()
// 		return nil, errors.New(msg)
// 	}
// 	defer tx.Rollback()
// 	repo := NewAuthRepository(tx)

// 	oldUser, err := repo.GetUserByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if oldUser.Type != 1 && oldUser.Type != 2 {
// 		msg := "用户类型错误"
// 		return nil, errors.New(msg)
// 	}
// 	byUser, err := repo.GetUserByID(byUserID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var byPriority int64
// 	byPriority = 0
// 	if byUser.RoleID != 0 {
// 		byRole, err := repo.GetRoleByID(byUser.RoleID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		byPriority = byRole.Priority
// 	}
// 	if oldUser.RoleID != 0 {
// 		targetRole, err := repo.GetRoleByID(oldUser.RoleID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if byPriority <= targetRole.Priority && userID != byUserID { //只能修改角色比自己优先级低的用户,或者用户自身
// 			msg := "你无法修改角色为" + targetRole.Name + "的用户"
// 			return nil, errors.New(msg)
// 		}
// 	}
// 	if info.RoleID != 0 {
// 		toRole, err := repo.GetRoleByID(info.RoleID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if byPriority < toRole.Priority { //只能将目标修改为和自己同级的角色
// 			msg := "你无法将目标角色改为:" + toRole.Name
// 			return nil, errors.New(msg)
// 		}
// 		oldUser.RoleID = info.RoleID
// 	}
// 	if info.PositionID != 0 {
// 		oldUser.PositionID = info.PositionID
// 	}
// 	if info.Name != "" {
// 		oldUser.Name = info.Name
// 	}
// 	if info.Email != "" {
// 		oldUser.Email = info.Email
// 	}
// 	if info.Gender != "" {
// 		oldUser.Gender = info.Gender
// 	}
// 	if info.Birthday != "" {
// 		oldUser.Birthday = info.Birthday
// 	}
// 	if info.Phone != "" {
// 		oldUser.Phone = info.Phone
// 	}
// 	if info.Address != "" {
// 		oldUser.Address = info.Address
// 	}
// 	if info.Status != 0 {
// 		if oldUser.ID != byUserID { //不能自己更新自己的状态
// 			oldUser.Status = info.Status
// 		} else if info.Status == 3 {
// 			oldUser.Status = info.Status
// 		}
// 	}
// 	if oldUser.Name == "" && oldUser.Status == 1 {
// 		msg := "必须有姓名才能启用用户"
// 		return nil, errors.New(msg)
// 	}
// 	err = repo.UpdateUser(userID, *oldUser, (*byUser).Name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := repo.GetUserByID(userID)
// 	tx.Commit()
// 	return user, err
// }

// func (s *authService) GetUserByID(id int64, organizationID int64) (*User, error) {
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	user, err := query.GetUserByID(id, organizationID)
// 	return user, err
// }

// func (s *authService) GetUserList(filter UserFilter, organizationID int64) (int, *[]UserResponse, error) {
// 	if organizationID != 0 && organizationID != filter.OrganizationID {
// 		filter.OrganizationID = organizationID
// 	}
// 	db := database.InitMySQL()
// 	query := NewAuthQuery(db)
// 	count, err := query.GetUserCount(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	list, err := query.GetUserList(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	return count, list, err
// }

// func (s *authService) UpdatePassword(info PasswordUpdate) error {
// 	db := database.InitMySQL()

// 	query := NewAuthQuery(db)
// 	credential, err := query.GetUserCredential(info.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	if !checkPasswordHash(info.OldPassword, credential) {
// 		errMessage := "旧密码错误"
// 		return errors.New(errMessage)
// 	}
// 	tx, err := db.Begin()
// 	if err != nil {
// 		msg := "事务开启错误" + err.Error()
// 		return errors.New(msg)
// 	}
// 	defer tx.Rollback()
// 	hashed, err := hashPassword(info.NewPassword)
// 	if err != nil {
// 		msg := "密码加密错误" + err.Error()
// 		return errors.New(msg)
// 	}
// 	repo := NewAuthRepository(tx)
// 	err = repo.UpdatePassword(info.UserID, hashed, info.User)
// 	if err != nil {
// 		msg := "密码更新错误" + err.Error()
// 		return errors.New(msg)
// 	}
// 	tx.Commit()
// 	return nil
// }
