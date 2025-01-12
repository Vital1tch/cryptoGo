package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"cryptoGo/data"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Функция для генерации ключа из пароля с использованием PBKDF2
func GenerateKeyFromPassword(password string, salt []byte) ([]byte, error) {
	// Генерация ключа из пароля с помощью PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New) // 32 - размер ключа (256 бит)
	return key, nil
}

// Генерация случайной соли для создания уникальных ключей
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16) // 16 байт соли
	if _, err := rand.Read(salt); err != nil {
		log.Println("Ошибка при генерации соли:", err)
		return nil, err
	}
	return salt, nil
}

// Функция для сохранения соли в файл
func SaveSalt(salt []byte, saltPath string) error {
	return os.WriteFile(saltPath, salt, 0600) // Сохраняем соль
}

// Функция для загрузки соли из файла
func LoadSalt(saltPath string) ([]byte, error) {
	fileData, err := os.ReadFile(saltPath)
	if err != nil {
		log.Println("Ошибка загрузки соли:", err)
		return nil, err
	}
	return fileData, nil
}

// Шифрование файла
func EncryptFileWithPassword(inputPath, outputPath, password string) error {
	salt, err := GenerateSalt()
	if err != nil {
		log.Println("Ошибка генерации соли:", err)
		return err
	}

	key, err := GenerateKeyFromPassword(password, salt)
	if err != nil {
		log.Println("Ошибка генерации ключа:", err)
		return err
	}

	// Проверяем, существует ли файл
	ext := filepath.Ext(inputPath)
	for {
		if _, err := os.Stat(outputPath); err == nil {
			log.Printf("Файл %s уже существует. Генерируется новое имя файла и соль.\n", outputPath)
			timestamp := time.Now().Format("20060102_150405")
			outputPath = fmt.Sprintf("%s_%s.enc", strings.TrimSuffix(outputPath, ext), timestamp)

			salt, err = GenerateSalt()
			if err != nil {
				log.Println("Ошибка генерации новой соли:", err)
				return err
			}

			key, err = GenerateKeyFromPassword(password, salt)
			if err != nil {
				log.Println("Ошибка генерации нового ключа:", err)
				return err
			}
		} else {
			break
		}
	}

	saltPath := outputPath + ".salt"
	err = SaveSalt(salt, saltPath)
	if err != nil {
		log.Println("Ошибка сохранения соли:", err)
		return err
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer inputFile.Close()

	fileData, err := io.ReadAll(inputFile)
	if err != nil {
		log.Println("Ошибка при чтении файла:", err)
		return err
	}

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

	// Добавляем оригинальное расширение в начало данных
	metadata := []byte(ext + "\n")
	dataWithMetadata := append(metadata, fileData...)

	ciphertext := aesGCM.Seal(nonce, nonce, dataWithMetadata, nil)

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

	// После успешного шифрования добавляем пароль в менеджер паролей
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	data.AddPassword(filepath.Base(outputPath), password, timestamp)

	log.Printf("Пароль для файла %s добавлен в менеджер паролей.\n", outputPath)

	log.Printf("Файл успешно зашифрован: %s, соль сохранена: %s\n", outputPath, saltPath)
	return nil
}

// Расшифровка файла
func DecryptFileWithPassword(inputPath, outputPath, password string) error {
	saltPath := inputPath + ".salt"
	salt, err := LoadSalt(saltPath)
	if err != nil {
		log.Println("Ошибка загрузки соли:", err)
		return err
	}

	key, err := GenerateKeyFromPassword(password, salt)
	if err != nil {
		log.Println("Ошибка генерации ключа:", err)
		return err
	}

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

	if len(ciphertext) < 12 {
		log.Println("Ошибка: недостаточная длина зашифрованных данных для извлечения nonce.")
		return fmt.Errorf("недостаточная длина зашифрованных данных")
	}

	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

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

	// Извлекаем оригинальное расширение из данных
	parts := strings.SplitN(string(plaintext), "\n", 2)
	if len(parts) < 2 {
		log.Println("Ошибка: не удалось извлечь метаданные из расшифрованных данных.")
		return fmt.Errorf("ошибка извлечения метаданных")
	}
	ext := parts[0]
	fileData := []byte(parts[1])

	// Указываем путь для сохранения расшифрованного файла
	filename := strings.TrimSuffix(filepath.Base(inputPath), ".enc")

	// Проверяем, есть ли уже расширение в имени файла
	if !strings.HasSuffix(filename, ext) {
		outputPath = filepath.Join("./decrypted", filename+ext)
	} else {
		outputPath = filepath.Join("./decrypted", filename)
	}

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
