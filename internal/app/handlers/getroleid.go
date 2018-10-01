package handlers

import (
	"net"

	"github.com/lheinrichde/golib/pkg/db"
	"github.com/lheinrichde/golib/pkg/handler"
)

// GetRoleID return role id
type GetRoleID int

// Handle connection
func (h GetRoleID) Handle(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// query database for role id
	var roleID int
	err = db.DB.QueryRow(`SELECT roles.roleID FROM roles
	INNER JOIN users ON users.role = roles.roleID
	WHERE users.username = $1;`, username).Scan(&roleID)
	if err != nil {
		return err
	}

	// set permissions and respond
	response := map[string]interface{}{"roleID": roleID}
	handler.Write(conn, response)

	return nil
}
