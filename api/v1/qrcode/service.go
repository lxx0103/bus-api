package qrcode

import (
	"bus-api/api/v1/user"
	"bus-api/core/database"
	"errors"
	"time"

	"github.com/google/uuid"
)

type qrcodeService struct {
}

func NewQrcodeService() *qrcodeService {
	return &qrcodeService{}
}

func (s *qrcodeService) NewWxQrcode(id int64) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewQrcodeRepository(tx)
	userRepo := user.NewUserRepository(tx)
	byUser, err := userRepo.GetWxUserByID(id)
	if err != nil {
		msg := "获取当前用户失败"
		return nil, errors.New(msg)
	}
	if byUser.Status != 2 { //已绑定
		msg := "用户状态错误"
		return nil, errors.New(msg)
	}
	count, err := repo.CheckQrCodePeriod(byUser.ID, time.Now().Add(time.Duration(-900*time.Second)))
	if err != nil {
		msg := "获取扫码记录失败"
		return nil, errors.New(msg)
	}
	if count {
		msg := "15分钟内已成功扫码过，无法再次生成"
		return nil, errors.New(msg)
	}
	err = repo.DeleteUserQrcode(id, byUser.Name)
	if err != nil {
		msg := "清除过时二维码失败"
		return nil, errors.New(msg)
	}
	var newQrcode Qrcode
	newQrcode.Code = uuid.New().String()
	newQrcode.UserID = id
	newQrcode.UserName = byUser.Name
	newQrcode.ExpireTime = time.Now().Add(time.Second * 30).Format("2006-01-02 15:04:05")
	newQrcode.Status = 1
	newQrcode.Created = time.Now()
	newQrcode.CreatedBy = byUser.Name
	newQrcode.Updated = time.Now()
	newQrcode.UpdatedBy = byUser.Name
	err = repo.CreateQrcode(newQrcode)
	if err != nil {
		msg := "创建二维码失败"
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &newQrcode.Code, nil
}

func (s *qrcodeService) ScanQrcode(info ScanQrcodeNew) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewQrcodeRepository(tx)
	userRepo := user.NewUserRepository(tx)
	byUser, err := userRepo.GetWxUserByID(info.UserID)
	if err != nil {
		msg := "获取当前用户失败"
		return errors.New(msg)
	}
	if byUser.Status != 2 {
		msg := "用户状态错误"
		return errors.New(msg)
	}
	qrcode, err := repo.GetQrcodeByCode(info.Code)
	if err != nil {
		msg := "二维码不存在或已过期"
		return errors.New(msg)
	}
	err = repo.ScanQrcode(qrcode.ID, byUser.Name)
	if err != nil {
		msg := "扫描二维码失败"
		return errors.New(msg)
	}
	var newHistory ScanHistory
	newHistory.Code = qrcode.Code
	newHistory.User = qrcode.UserName
	newHistory.ByUser = byUser.Name
	newHistory.UserID = qrcode.UserID
	newHistory.ByUserID = byUser.ID
	newHistory.ScanTime = time.Now().Format("2006-01-02 15:04:05")
	newHistory.Status = 1
	newHistory.Created = time.Now()
	newHistory.CreatedBy = byUser.Name
	newHistory.Updated = time.Now()
	newHistory.UpdatedBy = byUser.Name
	err = repo.CreateHistory(newHistory)
	if err != nil {
		msg := "创建扫码历史失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}

func (s *qrcodeService) GetHistoryList(filter HistoryFilter) (int, *[]ScanHistoryResponse, error) {
	db := database.RDB()
	query := NewQrcodeQuery(db)
	count, err := query.GetHistoryCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetHistoryList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}
