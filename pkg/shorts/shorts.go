package shorts

import (
	"log"

	"github.com/satori/go.uuid"
)

// GenerateUUID return UUID as string
func GenerateUUID() string {
	return uuid.NewV4().String()
}

// Check log error
func Check(err error) {
	if err != nil {
		log.Println(err)
	}
}
