package gui

import (
	"time"
)

type FileInfo struct {
	FilePath  string
	Operation string
	Timestamp string
}

// Список для хранения информации о файлах
var fileList []FileInfo

func AddFileToList(filePath, operation string) { // Функция для добавления информации о файле в список
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	fileList = append(fileList, FileInfo{
		FilePath:  filePath,
		Operation: operation,
		Timestamp: timestamp,
	})
}

func GetFileList() []FileInfo { // Функция для получения списка файлов
	return fileList
}
