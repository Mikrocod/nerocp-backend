package handlers

import (
	"database/sql"
	"net"

	"lheinrich.de/nerocp-backend/internal/pkg/db"
	"lheinrich.de/nerocp-backend/pkg/handler"
	"lheinrich.de/nerocp-backend/pkg/shorts"
)

// GetUsers return permissions
type GetUsers int

// Handle connection
func (h GetUsers) Handle(conn net.Conn, request map[string]interface{}, username string) {
	// check has permission
	if !handler.HasPermission(username, "page.userList") {
		handler.Error(conn, 403)
		return
	}

	// define variables
	role := handler.GetInt(request, "roleID")
	var rows *sql.Rows
	var err error

	// query database for users and check for error
	if role == 0 {
		// query for all users
		rows, err = db.DB.Query(`SELECT username, role FROM users;`)
	} else {
		// query for users with specific role
		rows, err = db.DB.Query(`SELECT username, role FROM users WHERE role = $1;`, role)
	}
	shorts.Check(err)

	// loop through rows
	users := []map[string]interface{}{}
	for rows.Next() {
		// define variables and scan
		var user string
		var role int
		rows.Scan(&user, &role)

		// add user to slice
		userItem := map[string]interface{}{}
		userItem["username"] = user
		userItem["roleID"] = role
		users = append(users, userItem)
	}

	// set users and respond
	response := map[string]interface{}{"users": users}
	handler.Write(conn, response)
}
