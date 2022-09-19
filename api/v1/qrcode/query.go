package qrcode

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type qrcodeQuery struct {
	conn *sqlx.DB
}

func NewQrcodeQuery(connection *sqlx.DB) *qrcodeQuery {
	return &qrcodeQuery{
		conn: connection,
	}
}

func (r *qrcodeQuery) GetHistoryCount(filter HistoryFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "user_id = ?"), append(args, v)
	}
	if v := filter.ByUserID; v != 0 {
		where, args = append(where, "by_user_id = ?"), append(args, v)
	}
	if v := filter.UserName; v != "" {
		where, args = append(where, "user like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ByUserName; v != "" {
		where, args = append(where, "by_user like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ScanDateFrom; v != "" {
		where, args = append(where, "scan_time >= ?"), append(args, v)
	}
	if v := filter.ScanDateTo; v != "" {
		where, args = append(where, "scan_time <= ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM q_scan_historys
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *qrcodeQuery) GetHistoryList(filter HistoryFilter) (*[]ScanHistoryResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.UserID; v != 0 {
		where, args = append(where, "user_id = ?"), append(args, v)
	}
	if v := filter.ByUserID; v != 0 {
		where, args = append(where, "by_user_id = ?"), append(args, v)
	}
	if v := filter.UserName; v != "" {
		where, args = append(where, "user like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ByUserName; v != "" {
		where, args = append(where, "by_user like ?"), append(args, "%"+v+"%")
	}
	if v := filter.ScanDateFrom; v != "" {
		where, args = append(where, "scan_time >= ?"), append(args, v)
	}
	if v := filter.ScanDateTo; v != "" {
		where, args = append(where, "scan_time <= ?"), append(args, v+" 23:59:59")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var wxQrcodes []ScanHistoryResponse
	err := r.conn.Select(&wxQrcodes, `
		SELECT id, code, user, user_id, by_user, by_user_id, scan_time, status
		FROM q_scan_historys
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &wxQrcodes, err
}
