package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // Генерация 256-битного ключа
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func SaveKey(key []byte, keyPath string) error {
	return os.WriteFile(keyPath, []byte(hex.EncodeToString(key)), 0600)
}

func LoadKey(keyPath string) ([]byte, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
