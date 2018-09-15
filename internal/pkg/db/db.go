package db

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"lheinrich.de/nerocp-backend/pkg/config"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// DB postgresql database
	DB *sql.DB
)

// Connect to postgresql database
func Connect() {
	// define variables
	var err error
	var trash string

	// connect
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

	// check if a role exists
	err = DB.QueryRow("SELECT roleID FROM roles;").Scan(&trash)
	if err == sql.ErrNoRows {
		// create default role
		_, err = DB.Exec("INSERT INTO roles (roleID, roleName) VALUES (1, 'admin')")
		shorts.Check(err)
	} else if err != nil {
		// print error
		shorts.Check(err)
	}

	// check if a user exists
	err = DB.QueryRow("SELECT username FROM users;").Scan(&trash)
	if err == sql.ErrNoRows {
		// generate bcrypt hash
		var passwordHash []byte
		passwordHash, err = bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost+1)
		shorts.Check(err)

		// create default user
		_, err = DB.Exec("INSERT INTO users (username, passwordHash, role) VALUES ('admin', $1, 1)", string(passwordHash))
		shorts.Check(err)
	} else if err != nil {
		// print error
		shorts.Check(err)
	}

}
