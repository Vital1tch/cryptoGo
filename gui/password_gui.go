package gui

import (
	"cryptoGo/data"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Открытие менеджера паролей с таблицей
func ShowPasswordManager(myWindow fyne.Window) {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль менеджера")

	dialog.ShowForm("Введите пароль менеджера", "ОК", "Отмена",
		[]*widget.FormItem{
			widget.NewFormItem("Пароль", passwordEntry),
		},
		func(confirm bool) {
			if confirm {
				if !data.CheckManagerPassword(passwordEntry.Text) {
					dialog.ShowError(errors.New("Неверный пароль"), myWindow)
					return
				}

				// Создаём заголовок таблицы
				header := container.NewGridWithColumns(3,
					widget.NewLabelWithStyle("Файл", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabelWithStyle("Пароль", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabelWithStyle("Время", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				)

				// Создаём строки таблицы
				var rows []fyne.CanvasObject
				rows = append(rows, header)
				for _, info := range data.PasswordList {
					// Обрезаем текст для каждого столбца
					filePath := info.FilePath
					if len(filePath) > 30 {
						filePath = filePath[:30] + "..."
					}

					password := info.Password
					if len(password) > 15 {
						password = password[:15] + "..."
					}

					row := container.NewGridWithColumns(3,
						widget.NewLabelWithStyle(filePath, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle(password, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle(info.Timestamp, fyne.TextAlignLeading, fyne.TextStyle{}),
					)
					rows = append(rows, row)
				}

				//// Оборачиваем таблицу в скроллинг
				table := container.NewVBox(rows...)
				//scrollableTable := container.NewScroll(table)

				// Показываем кастомный диалог с таблицей
				customDialog := dialog.NewCustom("Менеджер паролей", "Закрыть", table, myWindow)

				// Устанавливаем размер окна
				customDialog.Resize(fyne.NewSize(600, 300)) // Уменьшаем общий размер окна
				customDialog.Show()
			}
		}, myWindow,
	)
}
