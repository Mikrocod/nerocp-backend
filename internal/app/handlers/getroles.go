package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetRoles return permissions
type GetRoles int

// Handle connection
func (h GetRoles) Handle(conn net.Conn, request map[string]interface{}, username string) {
	// check has permission
	if !handler.HasPermission(username, "page.roleList") {
		handler.Error(conn, 403)
		return
	}

	// query database for roles
	rows, err := db.DB.Query(`SELECT roleID, roleName FROM roles;`)
	shorts.Check(err)

	// loop through rows
	roles := []map[string]interface{}{}
	for rows.Next() {
		// define variables and scan
		var roleID int
		var roleName string
		rows.Scan(&roleID, &roleName)

		// add user to slice
		role := map[string]interface{}{}
		role["roleID"] = roleID
		role["roleName"] = roleName
		roles = append(roles, role)
	}

	// set users and respond
	response := map[string]interface{}{"roles": roles}
	handler.Write(conn, response)
}
