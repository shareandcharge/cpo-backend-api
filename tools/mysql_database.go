package tools

import (
	"github.com/jmoiron/sqlx"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/motionwerkGmbH/msp-backend-api/tools"
)

var MDB *sqlx.DB

func MySQLConnect(dbName string) {

	//hardcoded, for now
	MDB, err := sqlx.Connect("mysql", "andy:hardpassword1@(18.197.172.83:3306)/"+dbName)
	tools.ErrorCheck(err, "mysql_database.go", true)

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
