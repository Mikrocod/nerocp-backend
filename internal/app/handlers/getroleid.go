package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetRoleID return role id
type GetRoleID int

// Handle connection
func (h GetRoleID) Handle(conn net.Conn, request map[string]interface{}) {
	// query database for role id
	var roleID int
	err := db.DB.QueryRow("SELECT roleID FROM users WHERE username = $1;", request["username"]).Scan(&roleID)
	shorts.Check(err)

	// set permissions and respond
	response := map[string]interface{}{"roleID": roleID}
	handler.Write(conn, response)
}
