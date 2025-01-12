package gui

import (
	"cryptoGo/crypto"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"path/filepath"
	"strings"
)

// Функция для отображения диалога ввода пароля
func ShowPasswordEntryDialog(title, message string, callback func(password string), parent fyne.Window) {
	// Поле ввода пароля
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	// Кнопки диалога
	dialogForm := dialog.NewForm(
		title,
		"OK",
		"Отмена",
		[]*widget.FormItem{
			widget.NewFormItem(message, passwordEntry),
		},
		func(confirm bool) {
			if confirm {
				callback(passwordEntry.Text)
			}
		},
		parent,
	)
	dialogForm.Resize(fyne.NewSize(400, 200))
	dialogForm.Show()
}

func StartApp() {
	// Загружаем список действий из файла
	LoadActionsFromFile()

	myApp := app.New()
	myWindow := myApp.NewWindow("CryptoApp")

	title := widget.NewLabel("Добро пожаловать в CryptoApp!")
	description := widget.NewLabel("Вы можете зашифровать и расшифровать свои файлы, используя только пароль!")
	action := widget.NewLabel("Выберите действие ниже:")
	actionsList := widget.NewLabel("Список взаимодействий с файлами:")
	en := widget.NewLabel(" ")

	// Создание таблицы для списка действий
	actionTable := widget.NewTable(
		func() (int, int) {
			return len(GetFileList()), 3 // Количество строк и столбцов
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Ячейки содержат текст
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			action := GetFileList()[id.Row]
			switch id.Col {
			case 0:
				cell.(*widget.Label).SetText(action.Timestamp)
			case 1:
				cell.(*widget.Label).SetText(action.Operation)
			case 2:
				cell.(*widget.Label).SetText(action.FilePath)
			}
		},
	)
	actionTable.SetColumnWidth(0, 150) // Время
	actionTable.SetColumnWidth(1, 150) // Действие
	actionTable.SetColumnWidth(2, 500) // Путь к файлу

	// Добавляем контейнер с прокруткой
	scrollableList := container.NewScroll(actionTable)
	scrollableList.SetMinSize(fyne.NewSize(600, 300))

	// Кнопка для открытия папки encrypted
	openEncryptedButton := widget.NewButton("Открыть папку с зашифрованными файлами", func() {
		OpenFolder("./encrypted")
	})

	// Кнопка для открытия папки decrypted
	openDecryptedButton := widget.NewButton("Открыть папку с расшифрованными файлами", func() {
		OpenFolder("./decrypted")
	})

	clearButton := widget.NewButton("Очистить список действий", func() {
		ClearActionList()
		actionTable.Refresh() // Обновляем таблицу
	})

	// Кнопки для шифрования и расшифровки
	encryptButton := widget.NewButton("Зашифровать файл", func() {
		ShowPasswordEntryDialog(
			"Введите пароль", "Пароль для шифрования",
			func(password string) {
				if password == "" {
					dialog.ShowError(errors.New("Пароль не может быть пустым"), myWindow)
					return
				}
				OpenFile(myWindow, func(inputPath string) {
					filename := filepath.Base(inputPath)
					outputPath := filepath.Join("./encrypted", filename+".enc")
					err := crypto.EncryptFileWithPassword(inputPath, outputPath, password)
					if err != nil {
						dialog.ShowError(err, myWindow)
						log.Printf("Ошибка шифрования файла: %s, ошибка: %v\n", inputPath, err)
						return
					}
					// Добавляем действие в список
					AddFileToList(inputPath, "Зашифрован")
					actionTable.Refresh()
					log.Printf("Файл успешно зашифрован: %s -> %s\n", inputPath, outputPath)
					dialog.ShowInformation("Успех", "Файл успешно зашифрован!", myWindow)
				})
			}, myWindow)
	})

	decryptButton := widget.NewButton("Расшифровать файл", func() {
		ShowPasswordEntryDialog(
			"Введите пароль", "Пароль для расшифровки",
			func(password string) {
				if password == "" {
					dialog.ShowError(errors.New("Пароль не может быть пустым"), myWindow)
					return
				}
				OpenEncFileFromDirectory(myWindow, "./encrypted", func(inputPath string) {
					filename := filepath.Base(inputPath)
					outputPath := filepath.Join("./decrypted", strings.TrimSuffix(filename, ".enc"))
					err := crypto.DecryptFileWithPassword(inputPath, outputPath, password)
					if err != nil {
						dialog.ShowError(err, myWindow)
						log.Printf("Ошибка расшифровки файла: %s, ошибка: %v\n", inputPath, err)
						return
					}
					// Добавляем действие в список
					AddFileToList(inputPath, "Расшифрован")
					actionTable.Refresh()
					log.Printf("Файл успешно расшифрован: %s -> %s\n", inputPath, outputPath)
					dialog.ShowInformation("Успех", "Файл успешно расшифрован!", myWindow)
				})
			}, myWindow)
	})

	// Главный экран
	myWindow.SetContent(container.NewVBox(
		title,
		description,
		action,
		container.NewHBox(encryptButton, decryptButton),
		en,
		container.NewHBox(openEncryptedButton, openDecryptedButton),
		en,
		actionsList,
		clearButton,
		scrollableList, // Добавляем прокручиваемую таблицу
	))
	myWindow.ShowAndRun()
}
