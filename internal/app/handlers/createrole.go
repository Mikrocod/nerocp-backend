package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/db"
	"github.com/lheinrichde/golib/pkg/handler"
)

// CreateRole create role
type CreateRole int

// Handle connection
func (h CreateRole) Handle(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// check if user has permission to create roles
	if !HasPermission(username, "page.roleList.create") {
		// no permission
		return errors.New("403")
	}

	// get role name
	roleName := handler.GetString(request, "roleName")

	// check if role name provided
	if roleName == "" {
		// role name is missing
		return errors.New("400")
	}

	// insert into database
	_, err = db.DB.Exec(`INSERT INTO roles (roleName) VALUES ($1);`, roleName)
	if err != nil {
		return err
	}

	// get role id
	var roleID string
	err = db.DB.QueryRow(`SELECT roleID FROM roles WHERE roleName = $1;`, roleName).Scan(&roleID)
	if err != nil {
		return err
	}

	// respond with success and role id
	err = handler.Write(conn, map[string]interface{}{"success": true, "roleID": roleID})
	if err != nil {
		return err
	}

	return nil
}
