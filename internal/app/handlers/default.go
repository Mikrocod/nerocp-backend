package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/pkg/handler"
)

// Default handle all requests
type Default int

// Handle connection
func (h Default) Handle(conn net.Conn, request map[string]interface{}, username string) {
	// respond with error 404
	response := map[string]interface{}{}
	response["error"] = 404
	handler.Write(conn, response)
}
