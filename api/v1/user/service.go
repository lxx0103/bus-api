package user

import (
	"bus-api/api/v1/auth"
	"bus-api/core/database"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
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
	identityValid := checkIDValid(info.Identity)
	if !identityValid {
		msg := "身份证验证失败"
		return errors.New(msg)
	}
	var wxUser auth.WxUser
	wxUser.OpenID = ""
	wxUser.Name = info.Name
	wxUser.School = info.School
	wxUser.Grade = info.Grade
	wxUser.Class = info.Class
	wxUser.Identity = info.Identity
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
	identityValid := checkIDValid(info.Identity)
	if !identityValid {
		msg := "身份证验证失败"
		return errors.New(msg)
	}
	var wxUser auth.WxUser
	wxUser.Name = info.Name
	wxUser.School = info.School
	wxUser.Grade = info.Grade
	wxUser.Class = info.Class
	wxUser.Identity = info.Identity
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
	wb, err := xlsx.OpenFile(path)
	if err != nil {
		msg := "Open excel Error!"
		return errors.New(err.Error() + msg)
	}
	sheetName := "Sheet1"
	sheet, ok := wb.Sheet[sheetName]
	if !ok {
		msg := "第一个Sheet必须名为： Sheet1"
		return errors.New(err.Error() + msg)
	}
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
	errmsg := ""
	sheet.ForEachRow(func(r *xlsx.Row) error {
		if r.GetCoordinate() == 0 {
			return nil
		}
		var wxUser auth.WxUser
		wxUser.Created = time.Now()
		wxUser.CreatedBy = byUser.Username
		wxUser.Updated = time.Now()
		wxUser.UpdatedBy = byUser.Username
		wxUser.Status = 1
		r.ForEachCell(func(c *xlsx.Cell) error {
			// cn, rn := c.GetCoordinates()
			cn, _ := c.GetCoordinates()
			switch cn {
			case 0:
				wxUser.School = c.Value
			case 1:
				wxUser.Grade = c.Value
			case 2:
				wxUser.Class = c.Value
			case 3:
				wxUser.Name = c.Value
			case 4:
				// if c.Value != "" {
				// 	identityValid := checkIDValid(c.Value)
				// 	if !identityValid {
				// 		msg := " 第" + strconv.Itoa(rn+1) + "行身份证验证失败"
				// 		errmsg += msg
				// 	}
				// 	wxUser.Identity = c.Value
				// } else {
				// 	wxUser.Identity = c.Value
				// }
				wxUser.Identity = c.Value
			}
			return err
		})
		if wxUser.Identity != "" {
			wxUsers = append(wxUsers, wxUser)
		}
		return err
	})
	if errmsg != "" {
		return errors.New(errmsg)
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
	oldWxUser, err := repo.GetWxUserByID(wxUserID)
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

func checkIDValid(id string) bool {
	// 身份证位数不对
	if len(id) != 15 && len(id) != 18 {
		return false
	}

	// 转大写
	id = strings.ToUpper(id)

	if len(id) == 18 {
		// 验证算法
		if !checkValidNo18(id) {
			fmt.Println(id, "身份证算法验证失败！")
			return false
		}

	} else {
		// 转18位
		id = idCard15To18(id)
	}

	// 生日验证
	if !checkBirthdayCode(id[6:14]) {
		fmt.Println(id, "生日验证失败！")
		return false
	}

	return true
}

//15位身份证转为18位
var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var validValue = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// 15位转18位
func idCard15To18(id15 string) string {
	nLen := len(id15)
	if nLen != 15 {
		return "身份证不是15位！"
	}
	id18 := make([]byte, 0)
	id18 = append(id18, id15[:6]...)
	id18 = append(id18, '1', '9')
	id18 = append(id18, id15[6:]...)

	sum := 0
	for i, v := range id18 {
		n, _ := strconv.Atoi(string(v))
		sum += n * weight[i]
	}
	mod := sum % 11
	id18 = append(id18, validValue[mod])
	return string(id18)
}

//18位身份证校验码
func checkValidNo18(id string) bool {
	//string -> []byte
	id18 := []byte(id)
	nSum := 0
	for i := 0; i < len(id18)-1; i++ {
		n, _ := strconv.Atoi(string(id18[i]))
		nSum += n * weight[i]
	}
	//mod得出18位身份证校验码
	mod := nSum % 11

	return validValue[mod] == id18[17]
}

// 验证生日
func checkBirthdayCode(birthday string) bool {
	year, _ := strconv.Atoi(birthday[:4])
	month, _ := strconv.Atoi(birthday[4:6])
	day, _ := strconv.Atoi(birthday[6:])

	curYear, curMonth, curDay := time.Now().Date()
	//出生日期大于现在的日期
	if year < 1900 || year > curYear || month <= 0 || month > 12 || day <= 0 || day > 31 {
		return false
	}

	if year == curYear {
		if month > int(curMonth) {
			return false
		} else if month == int(curMonth) && day > curDay {
			return false
		}
	}

	//出生日期在2月份
	if month == 2 {
		//闰年2月只有29号
		if isLeapYear(year) && day > 29 {
			return false
		} else if day > 28 { //非闰年2月只有28号
			return false
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 { //小月只有30号
		if day > 30 {
			return false
		}
	}

	return true
}

// 判断是否为闰年
func isLeapYear(year int) bool {
	if year <= 0 {
		return false
	}
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

func (s *userService) UnbindWxUser(wxUserID, byID int64) error {
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
	oldWxUser, err := repo.GetWxUserByID(wxUserID)
	if err != nil {
		msg := "获取微信用户失败"
		return errors.New(msg)
	}
	if oldWxUser.Status != 2 {
		msg := "此用户尚未绑定，无需解绑"
		return errors.New(msg)
	}
	var wxUser auth.WxUser
	wxUser.OpenID = ""
	wxUser.Status = 1
	wxUser.Updated = time.Now()
	wxUser.UpdatedBy = byUser.Username
	err = repo.UnbindWxUser(wxUserID, wxUser)
	if err != nil {
		msg := "解绑小程序用户失败"
		return errors.New(msg)
	}
	tx.Commit()
	return err
}
