package gui

import (
	"cryptoGo/crypto"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
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
			dialog.ShowError(err, window)
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
	absDir, err := filepath.Abs(directory)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	uri := storage.NewFileURI(absDir)
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
	fileDialog.SetLocation(lister)
	fileDialog.Show()
}

func HandleEncrypt(window fyne.Window, fileListData binding.StringList) {
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

		AddFileToList(outputPath, "Зашифрован")
		updateFileList(fileListData)
	})
}

func HandleDecrypt(window fyne.Window, fileListData binding.StringList) {
	OpenFileInDirectory(window, "./encrypted", func(inputPath string) {
		filename := filepath.Base(inputPath)
		baseFilename := strings.TrimSuffix(filename, ".enc")
		keyPath := filepath.Join("./keys", baseFilename+".key")

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

			AddFileToList(outputPath, "Расшифрован")
			updateFileList(fileListData)
		}, window)
	})
}

func updateFileList(fileListData binding.StringList) {
	var displayList []string
	for _, file := range fileList {
		displayList = append(displayList, truncatePath(file.FilePath, 30)+" - "+file.Operation+" - "+file.Timestamp)
	}
	fileListData.Set(displayList)
}
