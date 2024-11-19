package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

func DecryptFile(inputPath, outputPath, keyPath string) error {
	log.Println("Начало расшифровки файла:", inputPath)

	// Загрузка ключа
	key, err := LoadKey(keyPath)
	if err != nil {
		log.Println("Ошибка загрузки ключа:", err)
		return err
	}
	log.Println("Ключ успешно загружен из:", keyPath)

	// Чтение зашифрованного файла
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Println("Ошибка при открытии зашифрованного файла:", err)
		return err
	}
	defer inputFile.Close()

	ciphertext, err := io.ReadAll(inputFile)
	if err != nil {
		log.Println("Ошибка при чтении зашифрованного файла:", err)
		return err
	}
	log.Println("Зашифрованный файл успешно прочитан, размер данных:", len(ciphertext))

	if len(ciphertext) < 12 {
		err := errors.New("недостаточная длина зашифрованных данных")
		log.Println("Ошибка расшифровки:", err)
		return err
	}

	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	// Расшифровка данных
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Ошибка создания AES блока:", err)
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Ошибка создания GCM:", err)
		return err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Ошибка при расшифровке данных:", err)
		return err
	}

	// Извлечение расширения из данных
	data := string(plaintext)
	parts := strings.SplitN(data, "\n", 2)
	if len(parts) != 2 {
		err := errors.New("ошибка извлечения метаданных")
		log.Println(err)
		return err
	}
	extension := parts[0]
	fileData := []byte(parts[1])

	log.Println("Извлечено расширение файла:", extension)

	// Автоматическое добавление расширения к выходному файлу
	if !strings.HasSuffix(outputPath, extension) {
		outputPath += extension
	}
	log.Println("Итоговый путь выходного файла:", outputPath)

	// Запись расшифрованных данных
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Println("Ошибка создания выходного файла:", err)
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(fileData)
	if err != nil {
		log.Println("Ошибка записи расшифрованных данных:", err)
		return err
	}

	log.Println("Файл успешно расшифрован:", outputPath)
	return nil
}
