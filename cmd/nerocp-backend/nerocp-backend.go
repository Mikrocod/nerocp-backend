package main

import (
	"fmt"
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"

	"lheinrich.de/nerocp-backend/pkg/config"

	_ "github.com/lib/pq"
)

func main() {
	config.LoadConfig()
	db.Connect()
	fmt.Println("nerocp-backend (c) 2018 Lennart Heinrich")
	net.Listen("tcp", config.Get("server", "address"))
}
