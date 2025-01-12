package main

import (
	"cryptoGo/data"
	"cryptoGo/gui"
	"fyne.io/fyne/v2/app"
	"log"
	"os"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("CryptoGo")

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

	// Проверяем, установлен ли пароль менеджера
	if _, err := os.Stat(data.ManagerPasswordFile); os.IsNotExist(err) {
		log.Println("Пароль менеджера не найден. Пожалуйста, создайте его.")
		gui.ShowInitialPasswordDialog(myWindow)
	}

	// Создаём папки encrypted и decrypted, если их нет
	createFolderIfNotExists("./encrypted")
	createFolderIfNotExists("./decrypted")

	// Загружаем пароли
	data.LoadPasswordsFromFile() // Загружаем пароли из файла

	// Запускаем приложение
	gui.StartApp(myWindow)
}
