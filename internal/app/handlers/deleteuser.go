package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/handler"

	"github.com/lheinrichde/golib/pkg/db"
)

// DeleteUser function
func DeleteUser(conn net.Conn, request map[string]interface{}, username string) error {
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

	// respond with success
	err = handler.Write(conn, map[string]interface{}{"success": true})
	if err != nil {
		return err
	}

	return nil
}
