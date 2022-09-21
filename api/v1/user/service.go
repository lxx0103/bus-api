package user

import (
	"bufio"
	"bus-api/api/v1/auth"
	"bus-api/core/database"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"
)

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) GetWxUserByID(id int64) (*auth.WxUserResponse, error) {
	db := database.RDB()
	query := NewUserQuery(db)
	wxUser, err := query.GetWxUserByID(id)
	if err != nil {
		msg := "获取用户失败"
		return nil, errors.New(msg)
	}
	return wxUser, nil
}

func (s *userService) NewWxUser(info WxUserNew) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUserRepository(tx)
	authRepo := auth.NewAuthRepository(tx)
	byUser, err := authRepo.GetAdminUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	isConflict, err := repo.CheckWxIdentityConfict(0, info.Identity)
	if err != nil {
		msg := "身份证检验失败"
		return errors.New(msg)
	}
	if isConflict {
		msg := "身份证已存在"
		return errors.New(msg)
	}
	var wxUser auth.WxUser
	wxUser.OpenID = ""
	wxUser.Name = info.Name
	wxUser.Grade = info.Grade
	wxUser.Class = info.Class
	wxUser.Identity = info.Identity
	wxUser.Role = info.Role
	wxUser.ExpireDate = info.ExpireDate
	wxUser.Status = 1
	wxUser.Created = time.Now()
	wxUser.CreatedBy = byUser.Username
	wxUser.Updated = time.Now()
	wxUser.UpdatedBy = byUser.Username
	err = repo.CreateWxUser(wxUser)
	if err != nil {
		msg := "创建小程序用户失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}

func (s *userService) GetWxUserList(filter WxUserFilter) (int, *[]auth.WxUserResponse, error) {
	db := database.RDB()
	query := NewUserQuery(db)
	count, err := query.GetWxUserCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetWxUserList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *userService) UpdateWxUser(wxUserID int64, info WxUserNew) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUserRepository(tx)
	authRepo := auth.NewAuthRepository(tx)
	byUser, err := authRepo.GetAdminUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	isConflict, err := repo.CheckWxIdentityConfict(wxUserID, info.Identity)
	if err != nil {
		msg := "身份证检验失败"
		return errors.New(msg)
	}
	if isConflict {
		msg := "身份证已存在"
		return errors.New(msg)
	}
	var wxUser auth.WxUser
	wxUser.Name = info.Name
	wxUser.Grade = info.Grade
	wxUser.Class = info.Class
	wxUser.Identity = info.Identity
	wxUser.Role = info.Role
	wxUser.ExpireDate = info.ExpireDate
	wxUser.Updated = time.Now()
	wxUser.UpdatedBy = byUser.Username
	err = repo.UpdateWxUser(wxUserID, wxUser)
	if err != nil {
		msg := "更新小程序用户失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}

func (s *userService) DeleteWxUser(id, byID int64) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUserRepository(tx)
	authRepo := auth.NewAuthRepository(tx)
	byUser, err := authRepo.GetAdminUserByID(byID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	_, err = repo.GetWxUserByID(id)
	if err != nil {
		msg := "目标用户不存在"
		return errors.New(msg)
	}
	err = repo.DeleteWxUser(id, byUser.Username)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *userService) BatchUploadWxUser(path string, byID int64) error {
	csvFile, err := os.Open(path)
	if err != nil {
		msg := "打开文件错误"
		return errors.New(msg)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	authRepo := auth.NewAuthRepository(tx)
	byUser, err := authRepo.GetAdminUserByID(byID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	repo := NewUserRepository(tx)
	var wxUsers []auth.WxUser
	row := 0
	for {
		row = row + 1
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			msg := "读取文件错误"
			return errors.New(msg)
		}
		if row == 1 {
			continue
		}
		var wxUser auth.WxUser
		wxUser.Created = time.Now()
		wxUser.CreatedBy = byUser.Username
		wxUser.Updated = time.Now()
		wxUser.UpdatedBy = byUser.Username
		if line[0] == "" {
			msg := "第" + strconv.Itoa(row+1) + "行姓名为空"
			return errors.New(msg)
		} else {
			wxUser.Name = line[0]
		}
		if line[1] != "学生" && line[1] != "员工" {
			msg := "第" + strconv.Itoa(row+1) + "行角色错误，必须为学生或员工"
			return errors.New(msg)
		} else {
			wxUser.Role = line[1]
		}
		wxUser.Grade = line[2]
		wxUser.Class = line[3]
		if line[4] == "" {
			msg := "第" + strconv.Itoa(row+1) + "行身份证为空"
			return errors.New(msg)
		} else {
			wxUser.Identity = line[4]
		}
		t, err := time.Parse("2006-01-02", line[5])
		if err != nil {
			msg := "第" + strconv.Itoa(row+1) + "行日期格式错误"
			return errors.New(msg)
		} else {
			wxUser.ExpireDate = t.Format("2006-01-02")
		}
		if line[6] != "启用" && line[6] != "禁用" {
			msg := "第" + strconv.Itoa(row+1) + "行状态错误，必须为启用或禁用"
			return errors.New(msg)
		} else {
			if line[6] == "启用" {
				wxUser.Status = 1
			} else {
				wxUser.Status = 2
			}
		}
		wxUsers = append(wxUsers, wxUser)
	}
	err = repo.BatchCreateWxUser(wxUsers)
	if err != nil {
		msg := "批量创建用户失败" + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *userService) UpdateWxUserStatus(wxUserID int64, info WxUserStatusNew) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUserRepository(tx)
	authRepo := auth.NewAuthRepository(tx)
	byUser, err := authRepo.GetAdminUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	oldWxUser, err := repo.GetWxUserByID(info.UserID)
	if err != nil {
		msg := "获取微信用户失败"
		return errors.New(msg)
	}

	var wxUser auth.WxUser
	if info.Status == "active" {
		if oldWxUser.OpenID == "" {
			wxUser.Status = 1
		} else {
			wxUser.Status = 2
		}
	} else if info.Status == "deactive" {
		wxUser.Status = 3
	} else {
		msg := "状态错误"
		return errors.New(msg)
	}
	wxUser.Updated = time.Now()
	wxUser.UpdatedBy = byUser.Username
	err = repo.UpdateWxUserStatus(wxUserID, wxUser)
	if err != nil {
		msg := "更新小程序用户状态失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}
