package pkg

import (
	"encoding/json"
	"fmt"
)

// AnyPrint returns JSON-like string of any object with optional message
func AnyPrint(message string, i interface{}) string {
	JSON, err := json.MarshalIndent(i, "", "   ")
	if err != nil {
		fmt.Printf("AnyPrint: %s\n", err)
	}
	msg := "\n"
	if message != "" {
		msg = message + ":\n"
	}
	return msg + string(JSON)
}
