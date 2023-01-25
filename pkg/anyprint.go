package pkg

import (
	"encoding/json"
	"fmt"
)

// AnyPrint returns JSON-like string of any object with optional message
func AnyPrint(obj interface{}) string {
	JSON, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		fmt.Printf("AnyPrint: %s", err)
	}
	return string(JSON)
}
