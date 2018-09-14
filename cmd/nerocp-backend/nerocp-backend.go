package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"

	"lheinrich.de/nerocp-backend/internal/app/handlers"

	"lheinrich.de/nerocp-backend/pkg/handler"

	"lheinrich.de/nerocp-backend/pkg/module"

	"lheinrich.de/nerocp-backend/internal/pkg/db"

	"lheinrich.de/nerocp-backend/pkg/config"
	"lheinrich.de/nerocp-backend/pkg/shorts"

	_ "github.com/lib/pq"
)

// main function
func main() {
	// startup
	config.LoadConfig()
	setupLogging()
	db.Connect()
	startServer()
	registerHandlers()

	// external
	fmt.Println("nerocp-backend (c) 2018 Lennart Heinrich")
	module.LoadModules()

	// keep open
	<-make(chan bool)
}

// register known handlers
func registerHandlers() {
	handler.Add("getperms", handlers.GetPerms(0))
	handler.Add("getroleid", handlers.GetRoleID(0))
	handler.Add("getrolename", handlers.GetRoleName(0))
}

// setup logging
func setupLogging() {
	if config.Get("app", "logType") == "file" {
		// get log file
		logFile := time.Now().Format(config.Get("app", "logFile"))

		// split directory from filename and create them
		directory, _ := path.Split(logFile)
		os.MkdirAll(directory, os.ModePerm)

		// open file and check for error
		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		shorts.Check(err)

		// set file as output
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Println(err)
		}
	}
}

// start listener
func startServer() {
	// listen to address
	listener, errListen := net.Listen("tcp", config.Get("server", "address"))
	shorts.Check(errListen)

	// async
	go listen(listener)
}

// listener for incoming connections
func listen(listener net.Listener) {
	// add default handler
	handler.Add("default", handlers.Default(0))

	for {
		// accept connection
		conn, errConn := listener.Accept()
		shorts.Check(errConn)

		// handle connection
		go handleConnSafe(conn)
	}
}

// handle connection and close
func handleConnSafe(conn net.Conn) {
	handleConn(conn)
	conn.Close()
}

// handle connection
func handleConn(conn net.Conn) {
	// define variables
	request := handler.Read(conn)
	rawTyp, rawUsername, rawPassword := request["type"], request["username"], request["password"]

	// close connection if no type, username and password defined
	if rawTyp == nil || rawUsername == nil || rawPassword == nil {
		return
	}
	typ, username, password := rawTyp.(string), rawUsername.(string), rawPassword.(string)

	// verify login
	role := verifyLogin(username, password)
	response := map[string]interface{}{}

	// process verification
	if typ == "" || username == "" || password == "" || role == nil {
		// wrong login
		response["error"] = 403
		handler.Write(conn, response)
		return
	} else if typ == "login" {
		// respond with role id
		response["roleID"] = role
		handler.Write(conn, response)
		return
	}

	// handle with handler
	handler.Get(typ).Handle(conn, request)
}

// verify login
func verifyLogin(username, password string) *int {
	// query database for user role
	var role int
	row := db.DB.QueryRow("SELECT role FROM users WHERE username = $1 AND password = $2;", username, password).Scan(&role)

	// not exists or wrong login data
	if row == sql.ErrNoRows {
		return nil
	}

	// return role
	return &role
}
