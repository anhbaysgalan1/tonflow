package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeAES(data, pass string) (string, error) {
	src := []byte(data)
	key := sha256.Sum256([]byte(pass))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("new cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm error: %v", err)
	}

	nonce := key[len(key)-gcm.NonceSize():]

	dst := gcm.Seal(nil, nonce, src, nil)

	return hex.EncodeToString(dst), nil
}

func DecodeAES(data, pass string) (string, error) {
	src, err := hex.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("decode string error: %v", err)
	}

	key := sha256.Sum256([]byte(pass))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("new cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm error: %v", err)
	}

	nonce := key[len(key)-gcm.NonceSize():]

	decrypted, err := gcm.Open(nil, nonce, src, nil)
	if err != nil {
		return "", fmt.Errorf("decryption error: %v", err)
	}

	return string(decrypted), nil
}
