package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // Генерируем 256-битный ключ
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func SaveKey(key []byte, keyPath string) error {
	return os.WriteFile(keyPath, []byte(hex.EncodeToString(key)), 0600) // Сохраняем в hex-формате
}

func LoadKey(keyPath string) ([]byte, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}

func EncryptFile(inputPath, outputPath, keyPath string) error {
	log.Println("Начало шифрования файла:", inputPath)

	// Генерация ключа
	key, err := GenerateKey()
	if err != nil {
		log.Println("Ошибка генерации ключа:", err)
		return err
	}

	// Сохранение ключа
	err = SaveKey(key, keyPath)
	if err != nil {
		log.Println("Ошибка сохранения ключа:", err)
		return err
	}
	log.Println("Ключ успешно сохранен в:", keyPath)

	// Чтение содержимого файла
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer inputFile.Close()

	data, err := io.ReadAll(inputFile)
	if err != nil {
		log.Println("Ошибка при чтении файла:", err)
		return err
	}
	log.Println("Файл успешно прочитан, размер данных:", len(data))

	// Получение расширения файла
	extension := filepath.Ext(inputPath)
	if extension == "" {
		extension = ".bin" // По умолчанию, если расширение отсутствует
	}
	log.Println("Расширение файла:", extension)

	// Добавляем расширение файла к шифрованным данным
	metadata := []byte(extension + "\n") // Сохраняем расширение в первых байтах файла
	dataWithMetadata := append(metadata, data...)

	// Шифрование данных
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Ошибка создания AES блока:", err)
		return err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Ошибка генерации nonce:", err)
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Ошибка создания GCM:", err)
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, dataWithMetadata, nil)

	// Запись зашифрованных данных в новый файл
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Println("Ошибка создания выходного файла:", err)
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(ciphertext)
	if err != nil {
		log.Println("Ошибка записи зашифрованных данных:", err)
		return err
	}

	log.Println("Файл успешно зашифрован:", outputPath)
	return nil
}
