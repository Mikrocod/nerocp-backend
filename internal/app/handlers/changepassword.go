package handlers

import (
	"errors"
	"net"

	"golang.org/x/crypto/bcrypt"

	"github.com/lheinrichde/golib/pkg/db"

	"github.com/lheinrichde/golib/pkg/handler"
)

// ChangePassword function
func ChangePassword(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// get module name and check if provided
	newPassword, changeUsername := handler.GetString(request, "newPassword"), handler.GetString(request, "changeUsername")
	if newPassword == "" {
		return errors.New("400")
	}

	// check permission
	if changeUsername != "" && !HasPermission(username, "evortexcp.userList.editUser") {
		return errors.New("403")
	} else if changeUsername == "" {
		// set username for password change
		changeUsername = username
	}

	// hash password
	var passwordHash []byte
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost+1)
	if err != nil {
		return err
	}

	// update db
	_, err = db.DB.Exec("UPDATE users SET passwordHash = $1 WHERE username = $2;", string(passwordHash), changeUsername)
	if err != nil {
		return err
	}

	// respond
	handler.Write(conn, map[string]interface{}{"success": true})
	return nil
}
