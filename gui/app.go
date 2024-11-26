package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"os"
)

func truncatePath(path string, maxLength int) string {
	if len(path) > maxLength {
		return path[:maxLength] + "..."
	}
	return path
}

func StartApp() {
	myApp := app.New()
	myApp = app.NewWithID("com.crypto.app")
	myWindow := myApp.NewWindow("CryptoApp")

	title := widget.NewLabel("Добро пожаловать в CryptoApp!")

	encryptButton := widget.NewButton("Зашифровать", func() {
	})

	decryptButton := widget.NewButton("Расшифровать", func() {
	})

	// Связываем список файлов с отображением в UI
	fileListData := binding.NewStringList()

	encryptButton.OnTapped = func() {
		HandleEncrypt(myWindow, fileListData)
	}

	decryptButton.OnTapped = func() {
		HandleDecrypt(myWindow, fileListData)
	}

	// Список для отображения файлов
	filesList := widget.NewListWithData(fileListData,
		func() fyne.CanvasObject {
			return widget.NewLabel("File info")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	status := widget.NewLabel("")

	content := container.NewVBox(
		title,
		encryptButton,
		decryptButton,
		status,
	)
	myWindow.SetContent(container.NewBorder(content, nil, nil, nil, filesList))
	myWindow.Resize(fyne.NewSize(600, 400))
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
