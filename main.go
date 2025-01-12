package main

import (
	"cryptoGo/gui"
	"log"
	"os"
)

func main() {
	// Функция для проверки и создания директорий
	createFolderIfNotExists := func(folder string) {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			err := os.MkdirAll(folder, 0755)
			if err != nil {
				log.Fatalf("Не удалось создать папку %s: %v", folder, err)
			}
			log.Printf("Папка %s была успешно создана.", folder)
		} else {
			log.Printf("Папка %s уже существует.", folder)
		}
	}

	// Создаём папки encrypted и decrypted, если их нет
	createFolderIfNotExists("./encrypted")
	createFolderIfNotExists("./decrypted")

	// Запускаем приложение
	gui.StartApp()
}
