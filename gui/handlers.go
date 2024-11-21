package gui

import (
	"cryptoGo/crypto"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func OpenFileInDirectory(window fyne.Window, directory string, onFileSelected func(path string)) {
	absDir, err := filepath.Abs(directory) // Преобразуем в абсолютный путь
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	uri := storage.NewFileURI(absDir) // Используем абсолютный путь
	lister, err := storage.ListerForURI(uri)
	if err != nil {
		log.Printf("Ошибка создания ListableURI: %v", err)
		dialog.ShowError(err, window)
		return
	}

	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if reader == nil {
			dialog.ShowInformation("Информация", "Файл не выбран", window)
			return
		}
		defer reader.Close()
		onFileSelected(reader.URI().Path())
	}, window)

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".enc"}))
	log.Println("Абсолютный путь к директории:", absDir)
	log.Println("URI директории:", uri.String()) // Фильтруем только .enc файлы
	fileDialog.SetLocation(lister)               // Устанавливаем папку
	fileDialog.Show()
}

func HandleEncrypt(window fyne.Window) {
	OpenFile(window, func(inputPath string) {
		filename := filepath.Base(inputPath)
		keyPath := filepath.Join("./keys", filename+".key")
		outputPath := filepath.Join("./encrypted", filename+".enc")

		err := crypto.EncryptFile(inputPath, outputPath, keyPath)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		dialog.ShowInformation("Успех", "Файл зашифрован!\nКлюч сохранен в: "+keyPath+"\nЗашифрованный файл: "+outputPath, window)
	})
}

func HandleDecrypt(window fyne.Window) {
	OpenFileInDirectory(window, "./encrypted", func(inputPath string) {
		filename := filepath.Base(inputPath)
		// Удаляем расширение .enc из имени файла
		baseFilename := strings.TrimSuffix(filename, ".enc")
		keyPath := filepath.Join("./keys", baseFilename+".key")

		log.Println("Путь к зашифрованному файлу:", inputPath)
		log.Println("Имя файла:", filename)
		log.Println("Имя файла без .enc:", baseFilename)
		log.Println("Ожидаемый путь к ключу:", keyPath)

		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Ключ для файла не найден: "+keyPath), window)
			return
		}

		dialog.ShowEntryDialog("Сохранить файл", "Введите имя выходного файла", func(outputPath string) {
			err := crypto.DecryptFile(inputPath, outputPath, keyPath)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			dialog.ShowInformation("Успех", "Файл расшифрован!", window)
		}, window)
	})
}
