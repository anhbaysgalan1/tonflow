package pkg

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
)

// PrintAny returns JSON-like string of any object
func PrintAny(obj interface{}) string {
	JSON, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		fmt.Printf("PrintAny error: %s", err)
	}
	return string(JSON)
}

// ShortAddr returns first 4 and last 4 symbols of TON smartcontract address.
func ShortAddr(a string) string {
	return fmt.Sprintf("%s...%s", a[:4], a[len(a)-4:])
}

// FormatAmount parses float from a string and formats it in the American format with two decimal places
func FormatAmount(amount string) (string, error) {
	float, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("%.2f", float), nil
}
