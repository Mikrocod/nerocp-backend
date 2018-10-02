package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/module"

	"github.com/lheinrichde/golib/pkg/handler"
)

// DeinstallModule function
func DeinstallModule(conn net.Conn, request map[string]interface{}, username string) error {
	var err error

	// get module name and check if provided
	moduleName := handler.GetString(request, "moduleName")
	if moduleName == "" {
		return errors.New("400")
	}

	// check if user has permission
	if !HasPermission(username, "evortexcp.modules.deinstall") {
		return errors.New("403")
	}

	// get path and check if exists
	path := module.GetPath(moduleName)
	if path == "" {
		return errors.New("404")
	}

	// uninstall module
	err = module.RemoveModule(path)
	if err != nil {
		return err
	}

	// respond
	handler.Write(conn, map[string]interface{}{"success": true})
	return nil
}
