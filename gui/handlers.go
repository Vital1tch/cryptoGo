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

// Функция для обновления списка файлов в интерфейсе
func updateFileList(fileListData binding.StringList) {
	// Получаем список всех файлов в папке ./encrypted
	files, err := os.ReadDir("./encrypted")
	if err != nil {
		log.Printf("Ошибка при чтении папки ./encrypted: %v", err)
		return
	}

	// Создаем новый список строк для привязки
	var fileNames []string

	// Добавляем файлы в список
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	// Обновляем привязанный список данных
	fileListData.Set(fileNames)
}

//// Функция выбора файла для шифрования
//func OpenFile(window fyne.Window, onFileSelected func(path string)) {
//	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
//		if err != nil {
//			dialog.ShowError(err, window)
//			return
//		}
//		if reader == nil {
//			dialog.ShowInformation("Информация", "Файл не выбран", window)
//			return
//		}
//		defer reader.Close()
//		onFileSelected(reader.URI().Path()) // Передаем путь к выбранному файлу
//	}, window)
//}

// Функция для обработки шифрования
func HandleEncrypt(window fyne.Window, fileListData binding.StringList) {
	// Открытие диалога для выбора файла
	OpenFile(window, func(inputPath string) {
		// Выводим информацию о выбранном файле
		filename := filepath.Base(inputPath)
		keyPath := filepath.Join("./keys", filename+".key")
		outputPath := filepath.Join("./encrypted", filename+".enc")

		// Запрашиваем пароль для файла
		dialog.ShowEntryDialog("Введите пароль", "Пароль для шифрования файла", func(password string) {
			if password == "" {
				dialog.ShowError(errors.New("Пароль не может быть пустым"), window)
				return
			}

			// Шифруем файл
			err := crypto.EncryptFile(inputPath, outputPath, keyPath)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			// Сохраняем пароль для файла
			AddPasswordForFile(inputPath, password)

			// Показываем результат
			dialog.ShowInformation("Успех", "Файл зашифрован!\nКлюч сохранен в: "+keyPath+"\nЗашифрованный файл: "+outputPath, window)

			// Добавляем файл в список
			AddFileToList(outputPath, "Зашифрован")
			updateFileList(fileListData)
		}, window)
	})
}

// Функция для обработки расшифровки
func HandleDecrypt(window fyne.Window, fileListData binding.StringList) {
	OpenFileInDirectory(window, "./encrypted", func(inputPath string) {
		// Определяем имя файла и путь к ключу
		filename := filepath.Base(inputPath)
		baseFilename := strings.TrimSuffix(filename, ".enc")
		keyPath := filepath.Join("./keys", baseFilename+".key")

		// Проверяем наличие ключа
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Ключ для файла не найден: "+keyPath), window)
			return
		}

		// Запрашиваем имя выходного файла
		dialog.ShowEntryDialog("Сохранить файл", "Введите имя выходного файла", func(outputPath string) {
			err := crypto.DecryptFile(inputPath, outputPath, keyPath)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			// Показываем результат
			dialog.ShowInformation("Успех", "Файл расшифрован!", window)

			// Добавляем файл в список
			AddFileToList(outputPath, "Расшифрован")
			updateFileList(fileListData)
		}, window)
	})
}

// Функция для открытия файла в директории
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
