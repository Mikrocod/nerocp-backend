package module

import (
	"io/ioutil"
	"plugin"
	"strings"

	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// module list ID:Module
	modules = map[string]Module{}
)

// Module structure
type Module struct {
	Name    string
	Author  string
	Version string
}

// LoadModules load all plugins in modules
func LoadModules() {
	// list files
	files, errDir := ioutil.ReadDir("modules")
	shorts.Check(errDir)

	// loop through files and load modules
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			LoadModule("modules/" + file.Name())
		}
	}
}

// LoadModule load plugin from file
func LoadModule(path string) {
	// load plugin
	p, errOpen := plugin.Open(path)
	shorts.Check(errOpen)

	// load start function
	start, errStart := p.Lookup("Start")
	shorts.Check(errStart)

	// start module and add to map
	modules[path] = start.(func() Module)()
}

// RemoveModule disable plugin and remove
func RemoveModule(path string) {
	// load plugin
	p, errOpen := plugin.Open(path)
	shorts.Check(errOpen)

	// load stop function
	stop, errStop := p.Lookup("Stop")
	shorts.Check(errStop)

	// stop and remove module
	stop.(func())()
	delete(modules, path)
}
