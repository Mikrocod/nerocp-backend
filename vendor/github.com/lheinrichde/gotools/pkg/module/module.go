package module

import (
	"io/ioutil"
	"os"
	"plugin"
	"strings"
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
func LoadModules(dir string) error {
	var err error

	// create directory if not exists
	_, exists := os.Stat(dir)
	if os.IsNotExist(exists) {
		os.MkdirAll(dir, os.ModePerm)
	}

	// list files
	var files []os.FileInfo
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// loop through files and load modules
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			err = LoadModule(dir + "/" + file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadModule load plugin from file
func LoadModule(path string) error {
	var err error

	// load plugin
	var p *plugin.Plugin
	p, err = plugin.Open(path)
	if err != nil {
		return err
	}

	// load start function
	var start plugin.Symbol
	start, err = p.Lookup("Start")
	if err != nil {
		return err
	}

	// start module and add to map
	modules[path] = start.(func() Module)()

	return nil
}

// RemoveModule disable plugin and remove
func RemoveModule(path string) error {
	var err error

	// load plugin
	var p *plugin.Plugin
	p, err = plugin.Open(path)
	if err != nil {
		return err
	}

	// load stop function
	var stop plugin.Symbol
	stop, err = p.Lookup("Stop")
	if err != nil {
		return err
	}

	// stop and remove module
	stop.(func())()
	delete(modules, path)

	return nil
}
