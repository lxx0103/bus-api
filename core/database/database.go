package database

import (
	"fmt"
	"log"

	"bus-api/core/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var rdb *sqlx.DB
var wdb *sqlx.DB

func ConfigMysql() {
	host := config.ReadConfig("rdb.host")
	fmt.Println(host)
	user := config.ReadConfig("rdb.user")
	password := config.ReadConfig("rdb.password")
	port := config.ReadConfig("rdb.port")
	dbname := config.ReadConfig("rdb.dbname")
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&loc=Local"
	mysqldb, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error() + "rdb error")
	}
	mysqldb.SetMaxOpenConns(200)
	mysqldb.SetMaxIdleConns(10)
	rdb = mysqldb
	whost := config.ReadConfig("wdb.host")
	wuser := config.ReadConfig("wdb.user")
	wpassword := config.ReadConfig("wdb.password")
	wport := config.ReadConfig("wdb.port")
	wdbname := config.ReadConfig("wdb.dbname")
	wdsn := wuser + ":" + wpassword + "@tcp(" + whost + ":" + wport + ")/" + wdbname + "?parseTime=true&loc=Local"
	wmysqldb, err := sqlx.Connect("mysql", wdsn)
	if err != nil {
		log.Fatal(err.Error() + "wdb error")
	}
	wmysqldb.SetMaxOpenConns(200)
	wmysqldb.SetMaxIdleConns(10)
	wdb = wmysqldb
}

func WDB() *sqlx.DB {
	return wdb
}

func RDB() *sqlx.DB {
	return rdb
}
