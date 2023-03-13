package pkg

import (
	"bytes"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

func EncodeQR(s string) ([]byte, error) {
	encoder := qrcode.NewQRCodeWriter()

	img, err := encoder.Encode(s, gozxing.BarcodeFormat_QR_CODE, 512, 512, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode qr: %w", err)
	}

	buf := &bytes.Buffer{}

	err = png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}

	return buf.Bytes(), nil
}

func DecodeQR(r io.Reader) (string, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return "", fmt.Errorf("failed to decode jpeg: %w", err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to create bmp: %w", err)
	}

	decoder := qrcode.NewQRCodeReader()

	decoded, err := decoder.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decode qr: %w", err)
	}

	prefixes := []string{
		"ton://transfer/",
		"https://tonhub.com/transfer/",
		"https://test.tonhub.com/transfer/",
	}

	result := decoded.String()
	for _, p := range prefixes {
		if strings.HasPrefix(result, p) {
			result = strings.TrimPrefix(result, p)
		}
	}

	return result, nil
}
