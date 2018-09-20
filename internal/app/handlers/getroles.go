package handlers

import (
	"errors"
	"net"

	"lheinrich.de/lheinrich/golib/pkg/db"
	"lheinrich.de/lheinrich/golib/pkg/handler"
)

// GetRoles return permissions
type GetRoles int

// Handle connection
func (h GetRoles) Handle(conn net.Conn, request map[string]interface{}, username string) error {
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
