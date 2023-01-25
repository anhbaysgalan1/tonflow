package pkg

import (
	"crypto/aes"
	"encoding/hex"
)

func EncryptAES(key []byte, data string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(data))
	c.Encrypt(out, []byte(data))

	return hex.EncodeToString(out), nil
}

func DecryptAES(key []byte, ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return string(pt[:]), nil
}
