package handlers

import (
	"database/sql"
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetPerms return permissions
type GetPerms int

// Handle connection
func (h GetPerms) Handle(conn net.Conn, request map[string]interface{}, username string) {
	// define variables
	var rows *sql.Rows
	var err error
	roleID := handler.GetInt(request, "roleID")

	// check if role id is provided
	if roleID == 0 {
		// query database for permissions of requester
		rows, err = db.DB.Query(`SELECT permissions.permission FROM permissions
		INNER JOIN users ON users.role = permissions.role
		WHERE users.username = $1;`, username)
		shorts.Check(err)
	} else {
		// check if user has permission
		if handler.HasPermission(username, "page.roleList") {
			// query database for permissions of provided role
			rows, err = db.DB.Query(`SELECT permission FROM permissions
			WHERE role = $1;`, roleID)
			shorts.Check(err)
		}

		// no permission
		handler.Error(conn, 403)
		return
	}

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
