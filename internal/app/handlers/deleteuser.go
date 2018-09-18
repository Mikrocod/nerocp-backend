package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/gotools/pkg/handler"

	"github.com/lheinrichde/gotools/pkg/db"
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
		// delete from database
		_, err = db.DB.Exec(`DELETE FROM users WHERE username = $1;`, username)
		if err != nil {
			return err
		}
	} else {
		// check if user has permission
		if HasPermission(username, "page.userList.delete") {
			// delete from database by username
			_, err = db.DB.Exec(`DELETE FROM users WHERE username = $1;`, deleteUsername)
			if err != nil {
				return err
			}
		} else {
			// no permission
			return errors.New("403")
		}
	}

	return nil
}
