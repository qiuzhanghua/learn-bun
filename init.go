package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"

	"os"
)

// var sqldb *sql.DB
var db *bun.DB

func init() {
	godotenv.Load()
	activeDB := os.Getenv("ACTIVE_DB")
	//fmt.Println(activeDB)
	switch activeDB {
	case dialect.SQLite.String():
		initSqlite()
	case dialect.MySQL.String():
		initMysql()
	case dialect.PG.String():
		initPg()
	case dialect.MSSQL.String():
		initMssql()
	default:
		panic("Invalid ACTIVE_DB")
	}
}

func initSqlite() {
	sqlitePath, ok := os.LookupEnv("SQLITE_PATH")
	if !ok {
		sqlitePath = "memory"
	}
	dsn := "file::" + sqlitePath + ":?cache=shared"
	//fmt.Println(dsn)
	sqldb, err := sql.Open(sqliteshim.DriverName(), dsn)
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, sqlitedialect.New())
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func initMysql() {
	mysqlHost, ok := os.LookupEnv("MYSQL_HOST")
	if !ok {
		mysqlHost = "localhost"
	}
	mysqlPort, ok := os.LookupEnv("MYSQL_PORT")
	if !ok {
		mysqlPort = "3306"
	}
	mysqlUser, ok := os.LookupEnv("MYSQL_USER")
	if !ok {
		mysqlUser = "app"
	}
	mysqlPassword, ok := os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		mysqlPassword = "app"
	}
	mysqlDatabase, ok := os.LookupEnv("MYSQL_DATABASE")
	if !ok {
		mysqlDatabase = "app"
	}
	dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDatabase
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, mysqldialect.New())
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func initPg() {
	pgHost, ok := os.LookupEnv("PG_HOST")
	if !ok {
		pgHost = "localhost"
	}
	pgPort, ok := os.LookupEnv("PG_PORT")
	if !ok {
		pgPort = "5432"
	}
	pgUser, ok := os.LookupEnv("PG_USER")
	if !ok {
		pgUser = "app"
	}
	pgPassword, ok := os.LookupEnv("PG_PASSWORD")
	if !ok {
		pgPassword = "app"
	}
	pgDatabase, ok := os.LookupEnv("PG_DATABASE")
	if !ok {
		pgDatabase = "app"
	}
	dsn := "postgres://" + pgUser + ":" + pgPassword + "@" + pgHost + ":" + pgPort + "/" + pgDatabase + "?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

func initMssql() {
	mssqlHost, ok := os.LookupEnv("MSSQL_HOST")
	if !ok {
		mssqlHost = "localhost"
	}
	mssqlPort, ok := os.LookupEnv("MSSQL_PORT")
	if !ok {
		mssqlPort = "3306"
	}
	mssqlUser, ok := os.LookupEnv("MSSQL_USER")
	if !ok {
		mssqlUser = "app"
	}
	mssqlPassword, ok := os.LookupEnv("MSSQL_PASSWORD")
	if !ok {
		mssqlPassword = "app"
	}
	mssqlDatabase, ok := os.LookupEnv("MSSQL_DATABASE")
	if !ok {
		mssqlDatabase = "app"
	}
	dsn := "sqlserver://" + mssqlUser + ":" + mssqlPassword + "@" + mssqlHost + ":" + mssqlPort + "?datasource=" + mssqlDatabase
	sqldb, err := sql.Open("sqlserver", dsn)
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, mssqldialect.New())
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
