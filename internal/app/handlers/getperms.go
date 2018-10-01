package handlers

import (
	"database/sql"
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/db"
	"github.com/lheinrichde/golib/pkg/handler"
)

// GetPerms return permissions
type GetPerms int

// Handle connection
func (h GetPerms) Handle(conn net.Conn, request map[string]interface{}, username string) error {
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
		if err != nil {
			return err
		}
	} else {
		// check if user has permission
		if HasPermission(username, "page.roleList") {
			// query database for permissions of provided role
			rows, err = db.DB.Query(`SELECT permission FROM permissions
			WHERE role = $1;`, roleID)
			if err != nil {
				return err
			}
		} else {
			// no permission
			return errors.New("403")
		}
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

	return nil
}
