package db

import (
	"database/sql"

	"lheinrich.de/nerocp-backend/pkg/config"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// postgresql database
	db *sql.DB
)

// Connect to postgresql database
func Connect() {
	var err error
	db, err = sql.Open("postgres", "postgres://"+config.Get("postgresql", "username")+":"+config.Get("postgresql", "password")+"@"+config.Get("postgresql", "host")+"/"+config.Get("postgresql", "database")+"?sslmode="+config.Get("postgresql", "ssl"))
	shorts.Check(err)
}
