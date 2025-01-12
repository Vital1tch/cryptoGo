package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

// Открытие диалога для выбора файла из заданной директории
func OpenEncFileFromDirectory(window fyne.Window, directory string, onFileSelected func(path string)) {
	files, err := os.ReadDir(directory)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Создаём список файлов с расширением .enc
	fileOptions := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".enc") { // Проверяем расширение
			fileOptions = append(fileOptions, file.Name())
		}
	}

	// Если нет файлов с расширением .enc
	if len(fileOptions) == 0 {
		dialog.ShowInformation("Информация", "Нет файлов для расшифровки (.enc)", window)
		return
	}

	// Диалог выбора файла
	dialog.ShowCustom("Выберите файл", "Отмена",
		container.NewVBox( // Заменяем widget.NewVBox на container.NewVBox
			widget.NewLabel("Доступные файлы:"),
			widget.NewSelect(fileOptions, func(selected string) {
				if selected != "" {
					onFileSelected(filepath.Join(directory, selected)) // Возвращаем полный путь
				}
			}),
		), window)

}

// Открыть папку в проводнике
func OpenFolder(relativePath string) {
	// Получаем текущую рабочую директорию
	baseDir, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка получения текущей рабочей директории: %v\n", err)
		return
	}

	// Формируем полный путь к папке
	absPath := filepath.Join(baseDir, relativePath)

	log.Printf("Попытка открыть папку: %s (абсолютный путь: %s)\n", relativePath, absPath)

	if runtime.GOOS == "windows" {
		cmd := exec.Command("explorer", absPath)
		err := cmd.Start()
		if err != nil {
			log.Printf("Ошибка открытия папки %s: %v\n", absPath, err)
		} else {
			log.Printf("Папка успешно открыта: %s\n", absPath)
		}
	} else {
		log.Println("Открытие папок поддерживается только в Windows.")
	}
}
