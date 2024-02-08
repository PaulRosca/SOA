package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func connectDB() {
	var err error
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "test1234",
		Net:       "tcp",
		Addr:      "mysql.default.svc.cluster.local:3306",
		DBName:    "orders",
		ParseTime: true,
	}
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
}
