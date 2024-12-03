package gui

import (
	"cryptoGo/crypto"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
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
	myApp := app.New()
	myWindow := myApp.NewWindow("CryptoApp")

	title := widget.NewLabel("Добро пожаловать в CryptoApp!")

	encryptButton := widget.NewButton("Зашифровать", func() {
		ShowPasswordEntryDialog(
			"Введите пароль",
			"Пароль для шифрования",
			func(password string) {
				if password == "" {
					dialog.ShowError(errors.New("Пароль не может быть пустым"), myWindow)
					return
				}

				// Открыть диалог выбора файла
				OpenFile(myWindow, func(inputPath string) {
					filename := filepath.Base(inputPath)
					outputPath := filepath.Join("./encrypted", filename+".enc")
					err := crypto.EncryptFileWithPassword(inputPath, outputPath, password)
					if err != nil {
						dialog.ShowError(err, myWindow)
						return
					}
					dialog.ShowInformation("Успех", "Файл успешно зашифрован!", myWindow)
				})
			},
			myWindow,
		)
	})

	// Кнопка для расшифровки
	decryptButton := widget.NewButton("Расшифровать", func() {
		ShowPasswordEntryDialog(
			"Введите пароль",
			"Пароль для расшифровки",
			func(password string) {
				if password == "" {
					dialog.ShowError(errors.New("Пароль не может быть пустым"), myWindow)
					return
				}

				// Открыть диалог выбора зашифрованного файла
				OpenFile(myWindow, func(inputPath string) {
					filename := filepath.Base(inputPath)
					outputPath := filepath.Join("./decrypted", filename+".dec")

					// Вызов функции расшифровки с введенным паролем
					err := crypto.DecryptFileWithPassword(inputPath, outputPath, password)
					if err != nil {
						dialog.ShowError(err, myWindow)
						return
					}
					dialog.ShowInformation("Успех", "Файл успешно расшифрован!", myWindow)
				})
			},
			myWindow,
		)
	})

	myWindow.SetContent(container.NewVBox(title, encryptButton, decryptButton))
	myWindow.ShowAndRun()
}
