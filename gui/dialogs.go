package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
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
func OpenFileFromDirectory(window fyne.Window, directory string, onFileSelected func(path string)) {
	files, err := os.ReadDir(directory)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Создаем список файлов для выбора
	fileOptions := []string{}
	for _, file := range files {
		if !file.IsDir() {
			fileOptions = append(fileOptions, file.Name())
		}
	}

	// Если нет файлов
	if len(fileOptions) == 0 {
		dialog.ShowInformation("Информация", "Нет файлов для выбора", window)
		return
	}

	// Диалог выбора файла
	dialog.ShowCustom("Выберите файл", "Отмена",
		container.NewVBox( // Заменяем widget.NewVBox на container.NewVBox
			widget.NewLabel("Доступные файлы:"),
			widget.NewSelect(fileOptions, func(selected string) {
				if selected != "" {
					onFileSelected(directory + "/" + selected)
				}
			}),
		), window)
}
