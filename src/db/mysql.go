package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

//连接数据库
func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:xxx@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	//最大打开连接数
	db.SetMaxOpenConns(100)
	//闲置连接数
	db.SetMaxIdleConns(10)
	//连接生命周期（show variables like '%wait_timeout%'; ）
	db.SetConnMaxLifetime(time.Duration(28000) * time.Second)
	//建立实际连接
	if err = db.Ping(); err != nil {
		panic(err)
	}
}
