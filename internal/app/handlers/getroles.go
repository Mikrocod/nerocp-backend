package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/db"
	"github.com/lheinrichde/golib/pkg/handler"
)

// GetRoles function
func GetRoles(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// check has permission
	if !HasPermission(username, "page.roleList") {
		return errors.New("403")
	}

	// query database for roles
	rows, err := db.DB.Query(`SELECT roleID, roleName FROM roles;`)
	if err != nil {
		return err
	}
	defer rows.Close()

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

	return nil
}
