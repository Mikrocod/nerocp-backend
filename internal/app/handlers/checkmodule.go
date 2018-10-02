package handlers

import (
	"errors"
	"net"

	"github.com/lheinrichde/golib/pkg/handler"
	"github.com/lheinrichde/golib/pkg/module"
)

// CheckModule function
func CheckModule(conn net.Conn, request map[string]interface{}, username string) error {
	// get module name and check if provided
	moduleName := handler.GetString(request, "moduleName")
	if moduleName == "" {
		return errors.New("400")
	}

	// respond
	handler.Write(conn, map[string]interface{}{"installed": module.Exists(moduleName)})
	return nil
}
