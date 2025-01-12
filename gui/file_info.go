package gui

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type FileInfo struct {
	FilePath  string `json:"file_path"`
	Operation string `json:"operation"`
	Timestamp string `json:"timestamp"`
}

var fileList []FileInfo

const actionLogPath = "./actions.log"

func AddFileToList(filePath, operation string) {
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	action := FileInfo{
		FilePath:  filePath,
		Operation: operation,
		Timestamp: timestamp,
	}
	fileList = append(fileList, action)

	log.Printf("Добавлено действие: %s, файл: %s, время: %s\n", operation, filePath, timestamp)

	saveActionsToFile() // Сохраняем список действий в файл
}

func ClearActionList() {
	// Очищаем список действий в памяти
	fileList = []FileInfo{}

	// Перезаписываем файл actions.log пустым содержимым
	err := os.WriteFile(actionLogPath, []byte(""), 0600)
	if err != nil {
		log.Printf("Ошибка при очистке файла действий: %v\n", err)
		return
	}

	log.Println("Список действий и файл actions.log успешно очищены.")
}

// Получение списка действий
func GetFileList() []FileInfo {
	return fileList
}

func saveActionsToFile() {
	log.Println("Сохранение списка действий в файл...")
	file, err := os.Create(actionLogPath)
	if err != nil {
		log.Println("Ошибка при создании файла для действий:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(fileList)
	if err != nil {
		log.Println("Ошибка при сохранении действий:", err)
	} else {
		log.Println("Список действий успешно сохранён в файл.")
	}
}

func LoadActionsFromFile() {
	log.Println("Загрузка списка действий из файла...")
	file, err := os.Open(actionLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл действий не найден, будет создан новый.")
			return
		}
		log.Println("Ошибка при открытии файла действий:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fileList)
	if err != nil {
		log.Println("Ошибка при чтении файла действий:", err)
	} else {
		log.Printf("Загружено %d действий из файла.\n", len(fileList))
	}
}
