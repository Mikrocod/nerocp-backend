package db

import (
	"database/sql"

	"lheinrich.de/nerocp-backend/pkg/config"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// DB postgresql database
	DB *sql.DB
)

// Connect to postgresql database
func Connect() {
	// connect
	var err error
	DB, err = sql.Open("postgres", "postgres://"+config.Get("postgresql", "username")+":"+config.Get("postgresql", "password")+"@"+config.Get("postgresql", "host")+"/"+config.Get("postgresql", "database")+"?sslmode="+config.Get("postgresql", "ssl"))
	shorts.Check(err)

	// setup tables
	_, err = DB.Exec(`-- Setup
	-- Roles
	CREATE TABLE IF NOT EXISTS roles (roleID SERIAL UNIQUE, roleName VARCHAR(255) UNIQUE,
		PRIMARY KEY (roleID));

	-- Permissions
	CREATE TABLE IF NOT EXISTS permissions (role INT, permission VARCHAR(255),
		FOREIGN KEY (role) REFERENCES roles (roleID) ON DELETE CASCADE);

	-- Users
	CREATE TABLE IF NOT EXISTS users (username VARCHAR(255) UNIQUE, passwordHash VARCHAR(255), role INT,
		PRIMARY KEY(username), FOREIGN KEY (role) REFERENCES roles (roleID) ON DELETE CASCADE);`)
	shorts.Check(err)
}
