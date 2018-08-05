package tools

import (
	"github.com/jmoiron/sqlx"
	"time"
	_ "github.com/go-sql-driver/mysql"
)



var MDB *sqlx.DB

func MySQLConnect(dbName string) {
	MDB = sqlx.MustConnect("mysql", "andy:hardpassword1@(18.197.172.83:3306)/blockchain")

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
