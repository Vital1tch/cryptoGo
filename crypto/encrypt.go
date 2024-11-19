package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func EncryptFile(inputPath, outputPath, key string) error {
	log.Println("Начало шифрования файла:", inputPath)

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

	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		err := errors.New("длина ключа должна быть 16, 24 или 32 байта")
		log.Println("Ошибка ключа:", err)
		return err
	}

	// Шифрование данных
	block, err := aes.NewCipher([]byte(key))
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

	// Получение расширения файла
	extension := filepath.Ext(inputPath)
	if extension == "" {
		extension = ".bin" // По умолчанию, если у файла нет расширения
	}
	log.Println("Расширение файла:", extension)

	// Добавляем расширение файла к шифрованным данным
	metadata := []byte(extension + "\n") // Сохраняем расширение в начале
	ciphertext := aesGCM.Seal(nonce, nonce, append(metadata, data...), nil)

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
