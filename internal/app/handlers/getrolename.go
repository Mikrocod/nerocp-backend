package handlers

import (
	"net"

	"github.com/lheinrichde/gotools/pkg/db"
	"github.com/lheinrichde/gotools/pkg/handler"
)

// GetRoleName return role name
type GetRoleName int

// Handle connection
func (h GetRoleName) Handle(conn net.Conn, request map[string]interface{}, username string) error {
	// query database for role name
	var roleName string
	err := db.DB.QueryRow(`SELECT roles.roleName FROM roles
	INNER JOIN users ON users.role = roles.roleID
	WHERE username = $1;`, request["username"]).Scan(&roleName)
	if err != nil {
		return err
	}

	// set permissions and respond
	response := map[string]interface{}{"roleName": roleName}
	handler.Write(conn, response)

	return nil
}
