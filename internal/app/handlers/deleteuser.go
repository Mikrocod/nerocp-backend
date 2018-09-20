package handlers

import (
	"errors"
	"net"

	"lheinrich.de/lheinrich/golib/pkg/handler"

	"lheinrich.de/lheinrich/golib/pkg/db"
)

// DeleteUser delete user
type DeleteUser int

// Handle connection
func (h DeleteUser) Handle(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// get username to delete
	deleteUsername := handler.GetString(request, "deleteUsername")

	// check if username to delete is provided
	if deleteUsername == "" {
		// no username to delete provided
		return errors.New("400")
	}

	// check if user has permission
	if HasPermission(username, "page.userList.delete") || deleteUsername == username {
		// delete from database by username
		_, err = db.DB.Exec(`DELETE FROM users WHERE username = $1;`, deleteUsername)
		if err != nil {
			return err
		}
	} else {
		// no permission
		return errors.New("403")
	}

	return nil
}