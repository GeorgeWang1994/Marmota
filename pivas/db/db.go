package db

import (
	"database/sql"
	"log"
	"marmota/pivas/cc"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", cc.Config().Database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(cc.Config().MaxConns)
	DB.SetMaxIdleConns(cc.Config().MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
