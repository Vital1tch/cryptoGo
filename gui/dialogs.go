package gui

import (
	"cryptoGo/crypto"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// Открытие диалога для выбора файла
func OpenFile(window fyne.Window, onFileSelected func(path string)) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if reader == nil {
			dialog.ShowInformation("Информация", "Файл не выбран", window)
			return
		}
		defer reader.Close()
		onFileSelected(reader.URI().Path()) // Передаем путь к файлу
	}, window)
}

// Шифрование с паролем
func HandleEncryptWithPassword(window fyne.Window, password string, inputPath string, outputPath string) {
	err := crypto.EncryptFileWithPassword(inputPath, outputPath, password)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	dialog.ShowInformation("Успех", "Файл зашифрован успешно!", window)
}

// Расшифровка с паролем
func HandleDecryptWithPassword(window fyne.Window, password string, inputPath string, outputPath string) {
	err := crypto.DecryptFileWithPassword(inputPath, outputPath, password)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	dialog.ShowInformation("Успех", "Файл расшифрован успешно!", window)
}
