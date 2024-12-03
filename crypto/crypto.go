package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	data, err := os.ReadFile(saltPath)
	if err != nil {
		log.Println("Ошибка загрузки соли:", err)
		return nil, err
	}
	return data, nil
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

	data, err := io.ReadAll(inputFile)
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

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)

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

// Расшифровка файла с учетом пароля
func DecryptFileWithPassword(inputPath, outputPath, password string) error {
	// Создание директории для расшифровки, если она не существует
	outputDir := filepath.Dir(outputPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Println("Ошибка при создании директории для расшифрованного файла:", err)
			return err
		}
	}

	// Загружаем соль
	saltPath := inputPath + ".salt"
	salt, err := LoadSalt(saltPath)
	if err != nil {
		log.Println("Ошибка загрузки соли:", err)
		return err
	}

	// Генерация ключа из пароля и соли
	key, err := GenerateKeyFromPassword(password, salt)
	if err != nil {
		log.Println("Ошибка генерации ключа:", err)
		return err
	}

	// Открытие зашифрованного файла
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Println("Ошибка при открытии зашифрованного файла:", err)
		return err
	}
	defer inputFile.Close()

	// Чтение зашифрованного содержимого
	ciphertext, err := io.ReadAll(inputFile)
	if err != nil {
		log.Println("Ошибка при чтении зашифрованного файла:", err)
		return err
	}

	// Создание AES блока с использованием сгенерированного ключа
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Ошибка создания AES блока:", err)
		return err
	}

	// Создание GCM для работы с AES
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Ошибка создания GCM:", err)
		return err
	}

	// Извлечение nonce из зашифрованного файла
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	// Расшифровка данных
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Ошибка при расшифровке данных:", err)
		return err
	}

	// Получаем оригинальное расширение файла (до .enc)
	ext := filepath.Ext(inputPath)
	originalFilePath := strings.TrimSuffix(inputPath, ext) // Удаляем .enc

	// Создание выходного файла для расшифрованных данных
	outputFile, err := os.Create(originalFilePath)
	if err != nil {
		log.Println("Ошибка создания выходного файла:", err)
		return err
	}
	defer outputFile.Close()

	// Запись расшифрованных данных в файл
	_, err = outputFile.Write(plaintext)
	if err != nil {
		log.Println("Ошибка записи расшифрованных данных:", err)
		return err
	}

	log.Println("Файл успешно расшифрован:", originalFilePath)
	return nil
}
