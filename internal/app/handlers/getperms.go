package handlers

import (
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetPerms return permissions
type GetPerms int

// Handle connection
func (h GetPerms) Handle(conn net.Conn, request map[string]interface{}) {
	// close connection if no role id sent
	roleID, username, password := handler.GetInt(request, "roleID"), handler.GetString(request, "username"), handler.GetString(request, "password")
	if roleID == 0 {
		return
	}

	// query database for permissions
	rows, err := db.DB.Query(`SELECT permissions.permission FROM permissions
	INNER JOIN users ON permissions.role = users.role
	WHERE permissions.role = $1 AND username = $2 AND password = $3;`,
		roleID, username, password)
	shorts.Check(err)

	// loop through rows
	permissions := []string{}
	for rows.Next() {
		// scan and add to slice
		var permission string
		rows.Scan(&permission)
		permissions = append(permissions, permission)
	}

	// set permissions and respond
	response := map[string]interface{}{"permissions": permissions}
	handler.Write(conn, response)
}
