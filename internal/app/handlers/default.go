package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/pkg/handler"
)

// DefaultHandler handle all requests
type DefaultHandler int

// Handle connection
func (h DefaultHandler) Handle(conn net.Conn, request map[string]interface{}) {
	// respond with error 404
	response := map[string]interface{}{}
	response["error"] = 404
	handler.Write(conn, response)
}
