package handler

import (
	"encoding/json"
	"net"
	"strings"
)

var (
	// registered handlers
	handlersMap = map[string]Handler{}
)

// Handler interface for socket communication
type Handler interface {
	Handle(conn net.Conn, request map[string]interface{}, username string) error
}

// Get handler by name
func Get(name string) Handler {
	h, exists := handlersMap[strings.ToLower(name)]
	if !exists {
		// return default if not exists
		h = handlersMap["default"]
	}

	// return handler
	return h
}

// Add handler with name
func Add(name string, h Handler) {
	handlersMap[strings.ToLower(name)] = h
}

// Delete handler with name
func Delete(name string, h Handler) {
	delete(handlersMap, strings.ToLower(name))
}

// Read and unmarshal to map[string]interface{}
func Read(conn net.Conn) (map[string]interface{}, error) {
	var err error

	// read from connection
	requestBytes := make([]byte, 1024)
	var readLength int
	readLength, err = conn.Read(requestBytes)
	if err != nil {
		return nil, err
	}

	// unmarshal json
	var request map[string]interface{}
	err = json.Unmarshal(requestBytes[:readLength], &request)
	if err != nil {
		return nil, err
	}

	// return map
	return request, nil
}

// Write and marshal from map[string]interface{}
func Write(conn net.Conn, response map[string]interface{}) error {
	var err error

	// marshal to byte slice
	var responseBytes []byte
	responseBytes, err = json.Marshal(&response)
	if err != nil {
		return err
	}

	// write to connection
	_, err = conn.Write(responseBytes)
	if err != nil {
		return err
	}

	return nil
}

// GetString get string or empty string
func GetString(request map[string]interface{}, name string) string {
	// cast
	value, ok := request[name].(string)

	// return value
	if ok {
		return value
	}

	// return empty string
	return ""
}

// GetInt get int or 0
func GetInt(request map[string]interface{}, name string) int {
	// cast
	value, ok := request[name].(int)

	// check exists
	if ok {
		// return int
		return value
	}

	// cast float
	var float float64
	float, ok = request[name].(float64)

	// check float exists
	if ok {
		// return float as int
		return int(float)
	}

	// return 0
	return 0
}

// GetBool get string or false
func GetBool(request map[string]interface{}, name string) bool {
	// cast
	value, ok := request[name].(bool)

	// return value
	if ok {
		return value
	}

	// return false
	return false
}

// Error respond with error code
func Error(conn net.Conn, err error) {
	// define response and set error code
	response := map[string]interface{}{}
	response["error"] = err

	// Write
	Write(conn, response)
}
