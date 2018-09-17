package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/bcrypt"

	"github.com/lheinrichde/nerocp-backend/internal/app/handlers"

	"github.com/lheinrichde/gotools/pkg/config"
	"github.com/lheinrichde/gotools/pkg/crypter"
	"github.com/lheinrichde/gotools/pkg/db"
	"github.com/lheinrichde/gotools/pkg/handler"
	"github.com/lheinrichde/gotools/pkg/module"
	"github.com/lheinrichde/gotools/pkg/setup"

	_ "github.com/lib/pq"
)

// main function
func main() {
	var err error

	// config
	err = config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	// log
	if config.Get("app", "logType") == "file" {
		err = setup.LogToFile(config.Get("app", "logFile"))
		if err != nil {
			panic(err)
		}
	}

	// database
	err = setupDB()
	if err != nil {
		panic(err)
	}

	// server
	err = startServer()
	if err != nil {
		panic(err)
	}
	registerHandlers()

	// started
	fmt.Println("nerocp-backend (c) 2018 Lennart Heinrich")

	// modules
	err = module.LoadModules(config.Get("app", "modules"))
	if err != nil {
		panic(err)
	}

	// keep open
	<-make(chan bool)
}

// register known handlers
func registerHandlers() {
	handler.Add("getperms", handlers.GetPerms(0))
	handler.Add("getroleid", handlers.GetRoleID(0))
	handler.Add("getrolename", handlers.GetRoleName(0))
	handler.Add("getusers", handlers.GetUsers(0))
}

// connect to and setup database
func setupDB() error {
	var err error

	// unused sql query values
	var trash string

	// connect to databse
	db.Connect(config.Get("postgresql", "host"), config.Get("postgresql", "port"), config.Get("postgresql", "ssl"),
		config.Get("postgresql", "database"), config.Get("postgresql", "username"), config.Get("postgresql", "password"))

	// setup tables
	_, err = db.DB.Exec(`-- Setup
	-- Roles
	CREATE TABLE IF NOT EXISTS roles (roleID SERIAL UNIQUE, roleName VARCHAR(255) UNIQUE,
		PRIMARY KEY (roleID));

	-- Permissions
	CREATE TABLE IF NOT EXISTS permissions (role INT, permission VARCHAR(255),
		FOREIGN KEY (role) REFERENCES roles (roleID) ON DELETE CASCADE);

	-- Users
	CREATE TABLE IF NOT EXISTS users (username VARCHAR(255) UNIQUE, passwordHash VARCHAR(255), role INT,
		PRIMARY KEY(username), FOREIGN KEY (role) REFERENCES roles (roleID) ON DELETE CASCADE);`)
	if err != nil {
		return err
	}

	// check if a role exists
	err = db.DB.QueryRow("SELECT roleID FROM roles;").Scan(&trash)
	if err == sql.ErrNoRows {
		// create default role
		_, err = db.DB.Exec("INSERT INTO roles (roleID, roleName) VALUES (1, 'admin')")
		if err != nil {
			return err
		}
	} else if err != nil {
		// return error
		return err
	}

	// check if a user exists
	err = db.DB.QueryRow("SELECT username FROM users;").Scan(&trash)
	if err == sql.ErrNoRows {
		// generate bcrypt hash
		passwordHashSHA3 := crypter.Hash("nerocp_" + crypter.Hash("admin"))
		var passwordHash []byte
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(passwordHashSHA3), bcrypt.DefaultCost+1)
		if err != nil {
			return err
		}

		// create default user
		_, err = db.DB.Exec("INSERT INTO users (username, passwordHash, role) VALUES ('admin', $1, 1)", string(passwordHash))
		if err != nil {
			return err
		}
	} else if err != nil {
		// return error
		return err
	}

	return nil
}

// start listener
func startServer() error {
	var err error

	// listen to address
	var listener net.Listener
	listener, err = net.Listen("tcp", config.Get("server", "address"))
	if err != nil {
		return err
	}

	// async
	go listen(listener)

	return nil
}

// listener for incoming connections
func listen(listener net.Listener) {
	var err error

	// add default handler
	handler.Add("default", handlers.Default(0))

	for {
		// accept connection
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// handle connection
		go handleConnSafe(conn)
	}
}

// handle connection and close
func handleConnSafe(conn net.Conn) {
	var err error

	err = handleConn(conn)
	if err != nil {
		handler.Write(conn, map[string]interface{}{"error": err.Error()})
	}
	conn.Close()
}

// handle connection
func handleConn(conn net.Conn) error {
	var err error

	// define variables
	var request map[string]interface{}
	request, err = handler.Read(conn)
	if err != nil {
		return err
	}
	typ, username, password := handler.GetString(request, "type"), handler.GetString(request, "username"), handler.GetString(request, "password")

	// close connection if no type, username and password defined
	if typ == "" {
		return errors.New("404")
	} else if username == "" || password == "" {
		return errors.New("403")
	}

	// verify login
	role := verifyLogin(username, password)
	response := map[string]interface{}{}

	// process verification
	if role == nil {
		// wrong login
		return errors.New("403")
	} else if typ == "login" {
		// respond with role id
		response["roleID"] = *role
		handler.Write(conn, response)
	}

	// handle with handler
	handler.Get(typ).Handle(conn, request, username)

	return nil
}

// verify login
func verifyLogin(username, password string) *int {
	var err error

	// query database for user role
	var passwordHash string
	var role int
	err = db.DB.QueryRow("SELECT passwordHash, role FROM users WHERE username = $1", username).Scan(&passwordHash, &role)

	// not exists or wrong login data
	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		return nil
	}

	// return role
	return &role
}
