package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/db"
)

// Default function
func Default(conn net.Conn, request map[string]interface{}, username string) error {
	// respond with error 404
	return errors.New("404")
}

// HasPermission check user has permission
func HasPermission(username, permission string) bool {
	var err error
	var trash string

	// query
	err = db.DB.QueryRow(`SELECT permissions.permission FROM permissions
	INNER JOIN users ON users.role = permissions.role
	WHERE users.username = $1 AND (permissions.permission = $2 OR permissions.permission = '*');`,
		username, permission).Scan(&trash)

	// check if has permission
	if err == nil {
		// return true
		return true
	}

	return false
}
