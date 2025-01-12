package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"os"
)

type PasswordInfo struct {
	FilePath  string `json:"file_path"`
	Password  string `json:"password"`
	Timestamp string `json:"timestamp"`
}

const (
	PasswordFile        = "./passwords.json"
	ManagerPasswordFile = "./manager_password.hash"
	saltSize            = 16
)

var PasswordList []PasswordInfo

// Генерация ключа для шифрования на основе пароля
func generateKey(password string, salt []byte) ([]byte, error) {
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New) // 32 байта = 256 бит
	return key, nil
}

// Шифрование данных
func encryptData(data []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	key, err := generateKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}

// Расшифровка данных
func decryptData(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < saltSize {
		return nil, errors.New("invalid ciphertext")
	}

	salt := ciphertext[:saltSize]
	ciphertext = ciphertext[saltSize:]

	key, err := generateKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// AddPassword добавляет пароль в список и сохраняет в файл
func AddPassword(filePath, password, timestamp string) {
	log.Printf("Добавляем пароль для файла: %s", filePath)
	PasswordList = append(PasswordList, PasswordInfo{
		FilePath:  filePath,
		Password:  password,
		Timestamp: timestamp,
	})
	savePasswordsToFile()
}

// savePasswordsToFile сохраняет пароли в файл с шифрованием
func savePasswordsToFile() {
	data, err := json.Marshal(PasswordList)
	if err != nil {
		log.Println("Ошибка при сериализации паролей:", err)
		return
	}

	managerPassword := loadManagerPassword()
	if managerPassword == "" {
		log.Println("Пароль менеджера не установлен. Пароли не будут сохранены.")
		return
	}

	encryptedData, err := encryptData(data, managerPassword)
	if err != nil {
		log.Println("Ошибка при шифровании паролей:", err)
		return
	}

	err = os.WriteFile(PasswordFile, encryptedData, 0600)
	if err != nil {
		log.Println("Ошибка при записи файла паролей:", err)
	} else {
		log.Println("Пароли успешно сохранены (зашифрованы).")
	}
}

// LoadPasswordsFromFile загружает пароли из файла с расшифровкой
func LoadPasswordsFromFile() {
	log.Println("Загрузка паролей из файла...")
	encryptedData, err := os.ReadFile(PasswordFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл паролей не найден, будет создан новый.")
			return
		}
		log.Println("Ошибка при чтении файла паролей:", err)
		return
	}

	managerPassword := loadManagerPassword()
	if managerPassword == "" {
		log.Println("Пароль менеджера не установлен. Невозможно загрузить пароли.")
		return
	}

	data, err := decryptData(encryptedData, managerPassword)
	if err != nil {
		log.Println("Ошибка при расшифровке паролей:", err)
		return
	}

	err = json.Unmarshal(data, &PasswordList)
	if err != nil {
		log.Println("Ошибка при десериализации паролей:", err)
	} else {
		log.Printf("Загружено %d паролей из файла.\n", len(PasswordList))
	}
}

// SetManagerPassword устанавливает хэшированный пароль для менеджера паролей
func SetManagerPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = os.WriteFile(ManagerPasswordFile, hash, 0600)
	if err != nil {
		log.Println("Ошибка при сохранении пароля менеджера:", err)
		return err
	}
	log.Println("Пароль менеджера успешно установлен.")
	return nil
}

// CheckManagerPassword проверяет введённый пароль
func CheckManagerPassword(password string) bool {
	hash, err := os.ReadFile(ManagerPasswordFile)
	if err != nil {
		log.Println("Ошибка при чтении пароля менеджера:", err)
		return false
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

// ResetManagerPassword сбрасывает пароль менеджера и очищает данные о паролях
func ResetManagerPassword() error {
	err := os.Remove(ManagerPasswordFile)
	if err != nil && !os.IsNotExist(err) {
		log.Println("Ошибка при удалении файла пароля менеджера:", err)
		return err
	}

	err = os.Remove(PasswordFile)
	if err != nil && !os.IsNotExist(err) {
		log.Println("Ошибка при удалении файла с паролями:", err)
		return err
	}

	PasswordList = []PasswordInfo{}
	log.Println("Пароли успешно удалены из памяти и файла.")

	return nil
}

// Загружает пароль менеджера из хэша
func loadManagerPassword() string {
	hash, err := os.ReadFile(ManagerPasswordFile)
	if err != nil {
		log.Println("Пароль менеджера не найден.")
		return ""
	}
	return string(hash)
}
