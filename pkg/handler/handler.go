package handler

import (
	"encoding/json"
	"net"

	"lheinrich.de/nerocp-backend/pkg/shorts"
)

var (
	// registered handlers
	handlersMap = map[string]Handler{}
)

// Handler interface for socket communication
type Handler interface {
	Handle(conn net.Conn, request map[string]interface{})
}

// Get handler by name
func Get(name string) Handler {
	h, exists := handlersMap[name]
	if !exists {
		// return default if not exists
		h = handlersMap["default"]
	}
	return h
}

// Add handler with name
func Add(name string, h Handler) {
	handlersMap[name] = h
}

// Read and unmarshal to map[string]interface{}
func Read(conn net.Conn) map[string]interface{} {
	var err error

	// read from connection
	var requestBytes []byte
	_, err = conn.Read(requestBytes)
	shorts.Check(err)

	// unmarshal json
	var request map[string]interface{}
	err = json.Unmarshal(requestBytes, &request)
	shorts.Check(err)

	// return map
	return request
}

// Write and marshal from map[string]interface{}
func Write(conn net.Conn, response map[string]interface{}) {
	var err error

	// marshal to byte slice
	var responseBytes []byte
	responseBytes, err = json.Marshal(&response)
	shorts.Check(err)

	// write to connection
	_, err = conn.Write(responseBytes)
	shorts.Check(err)
}