package tools

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type CPO struct {
	CpoId       int    `db:"cpo_id"`
	PublicAddr  string `db:"public_addr"`
	Seed string `db:"seed"`
	Email       string `db:"email"`
	Password    string `db:"password"`
}

var DB *sqlx.DB

func Connect(dbName string) {
	DB = sqlx.MustConnect("sqlite3", dbName)

	//some benchmark should be done here
	DB.SetMaxOpenConns(300)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(10 * time.Second)
}
