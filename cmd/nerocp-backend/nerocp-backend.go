package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"plugin"
	"time"

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

	// external
	fmt.Println("nerocp-backend (c) 2018 Lennart Heinrich")
	module.LoadModules()

	// keep open
	<-make(chan bool)
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
	for {
		// accept connection
		conn, errConn := listener.Accept()
		shorts.Check(errConn)

		// handle connection
		go handleConn(conn)
	}
}

// handle connection
func handleConn(conn net.Conn) {
	p, _ := plugin.Open("")
	s, _ := p.Lookup("")
	s.(func())()
}
