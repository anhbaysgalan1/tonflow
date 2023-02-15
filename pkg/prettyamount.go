package pkg

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
)

func PrettyAmount(amount string) (string, error) {
	float, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	p := message.NewPrinter(language.AmericanEnglish)
	return p.Sprintf("%.2f", float), nil
}
