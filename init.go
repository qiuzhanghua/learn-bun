package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var sqldb *sql.DB
var db *bun.DB

func init() {

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, sqlitedialect.New())

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqldb, err = sql.Open("mysql", "app:app@/app")
	if err != nil {
		panic(err)
	}

	db = bun.NewDB(sqldb, mysqldialect.New())

	err = db.Ping()
	if err != nil {
		panic(err)
	}

}
