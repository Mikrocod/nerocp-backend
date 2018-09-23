package handlers

import (
	"database/sql"
	"errors"
	"net"

	"lheinrich.de/lheinrich/golib/pkg/db"
	"lheinrich.de/lheinrich/golib/pkg/handler"
)

// GetUsers return permissions
type GetUsers int

// Handle connection
func (h GetUsers) Handle(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// check has permission
	if !HasPermission(username, "page.userList") {
		return errors.New("403")
	}

	// define variables
	role := handler.GetInt(request, "roleID")
	var rows *sql.Rows

	// query database for users and check for error
	if role == 0 {
		// query for all users
		rows, err = db.DB.Query(`SELECT users.username, users.role, roles.roleName
		FROM users INNER JOIN roles
		ON roles.roleID = users.role;`)
	} else {
		// query for users with specific role
		rows, err = db.DB.Query(`SELECT users.username, users.role, roles.roleName
		FROM users INNER JOIN roles
		ON roles.roleID = users.role
		WHERE users.role = $1;`, role)
	}
	if err != nil {
		return err
	}

	// loop through rows
	users := []map[string]interface{}{}
	for rows.Next() {
		// define variables and scan
		var user string
		var roleID int
		var roleName string
		rows.Scan(&user, &roleID, &roleName)

		// add user to slice
		userItem := map[string]interface{}{}
		userItem["username"] = user
		userItem["roleID"] = roleID
		userItem["roleName"] = roleName
		users = append(users, userItem)
	}

	// set users and respond
	response := map[string]interface{}{"users": users}
	handler.Write(conn, response)

	return nil
}
