package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/db"
	"github.com/lheinrichde/golib/pkg/handler"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser functio
func CreateUser(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// check if user has permission to create users
	if !HasPermission(username, "page.userList.create") {
		// no permission
		return errors.New("403")
	}

	// get new username, password and role id
	newUsername, newPassword := handler.GetString(request, "newUsername"), handler.GetString(request, "newPassword")
	newRoleID := handler.GetInt(request, "newRoleID")

	// check if all data provided
	if newUsername == "" || newPassword == "" || newRoleID == 0 {
		// something is missing
		return errors.New("400")
	}

	// hash password with bcrypt
	var passwordHash []byte
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost+1)
	if err != nil {
		return err
	}

	// insert into database
	_, err = db.DB.Exec(`INSERT INTO users (username, passwordHash, role) VALUES ($1, $2, $3);`, newUsername, string(passwordHash), newRoleID)
	if err != nil {
		return err
	}

	// respond with success
	err = handler.Write(conn, map[string]interface{}{"success": true})
	if err != nil {
		return err
	}

	return nil
}
