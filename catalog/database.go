package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func connectDB() {
	var err error
	cfg := mysql.Config{
		User:   "root",
		Passwd: "test1234",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "catalog",
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
