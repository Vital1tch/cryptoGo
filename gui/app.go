package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

func StartApp() {
	myApp := app.New()
	myApp = app.NewWithID("com.crypto.app")
	myWindow := myApp.NewWindow("CryptoApp")

	title := widget.NewLabel("Добро пожаловать в CryptoApp!")

	encryptButton := widget.NewButton("Зашифровать", func() {
		HandleEncrypt(myWindow)
	})

	decryptButton := widget.NewButton("Расшифровать", func() {
		HandleDecrypt(myWindow)
	})

	status := widget.NewLabel("")

	content := container.NewVBox(
		title,
		encryptButton,
		decryptButton,
		status,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 200))
	myWindow.ShowAndRun()
}

func EnsureDirectories() {
	dirs := []string{"./keys", "./encrypted"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755) // Создаем папку с правами доступа
		}
	}
}
