package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/handler"

	"github.com/lheinrichde/golib/pkg/db"
)

// DeleteRole function
func DeleteRole(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// get role id to delete
	roleID := handler.GetInt(request, "roleID")

	// check if role id provided
	if roleID == 0 {
		// no role id provided
		return errors.New("400")
	}

	// check if user has permission
	if HasPermission(username, "page.roleList.delete") {
		// delete from database by role id
		_, err = db.DB.Exec(`DELETE FROM roles WHERE roleID = $1;`, roleID)
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
