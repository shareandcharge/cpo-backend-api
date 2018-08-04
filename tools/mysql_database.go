package tools

import (
	"github.com/jmoiron/sqlx"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

//topics:
//msptokenmint: "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"
//cdr: "0x584eba3e36b122adb57f1e2f4ed1eb8e2b1bbbe1dc6241374239e45e4520a47d"
//error: "0x57cf7a55e859b30b6bfeb9a7dd14411606106cb3e082f2cda387ec3b4b90be1c"
//transferToken: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

var MDB *sqlx.DB

func MySQLConnect(dbName string) {
	MDB = sqlx.MustConnect("mysql", "andy:hardpassword1@(18.197.172.83:3306)/blockchain")

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
