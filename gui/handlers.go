package gui

import (
	"cryptoGo/crypto"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func OpenFile(window fyne.Window, onFileSelected func(path string)) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window) // Показываем ошибку
			return
		}
		if reader == nil {
			dialog.ShowInformation("Информация", "Файл не выбран", window)
			return
		}
		defer reader.Close()
		onFileSelected(reader.URI().Path()) // Передаём путь к выбранному файлу
	}, window)
}

func HandleEncrypt(window fyne.Window) {
	OpenFile(window, func(inputPath string) {
		dialog.ShowEntryDialog("Сохранить файл", "Введите имя выходного файла", func(outputPath string) {
			key := "examplekey123456" // Ключ длиной 16 байт
			err := crypto.EncryptFile(inputPath, outputPath, key)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			dialog.ShowInformation("Успех", "Файл зашифрован!", window)
		}, window)
	})
}

func HandleDecrypt(window fyne.Window) {
	OpenFile(window, func(inputPath string) {
		dialog.ShowEntryDialog("Сохранить файл", "Введите имя выходного файла", func(outputPath string) {
			key := "examplekey123456" // Ключ для расшифровки
			err := crypto.DecryptFile(inputPath, outputPath, key)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			dialog.ShowInformation("Успех", "Файл расшифрован!", window)
		}, window)
	})
}
