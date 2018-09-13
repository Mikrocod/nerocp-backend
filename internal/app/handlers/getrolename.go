package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetRoleName return role name
type GetRoleName int

// Handle connection
func (h GetRoleName) Handle(conn net.Conn, request map[string]interface{}) {
	// query database for role name
	var roleName string
	err := db.DB.QueryRow("SELECT roleName FROM users WHERE username = $1;", request["username"]).Scan(&roleName)
	shorts.Check(err)

	// set permissions and respond
	response := map[string]interface{}{"roleName": roleName}
	handler.Write(conn, response)
}
