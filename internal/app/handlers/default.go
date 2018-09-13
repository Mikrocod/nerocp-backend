package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/pkg/handler"
)

// DefaultHandler handle all requests
type DefaultHandler int

// Handle connection
func (h DefaultHandler) Handle(conn net.Conn, request map[string]interface{}) {
	handler.Write(conn, map[string]interface{}{})
}
