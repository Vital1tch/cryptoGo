package gui

import (
	"cryptoGo/data"
	"fyne.io/fyne/v2/widget"
	"time"
)

func AddPasswordForFile(filePath, password string) {
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	data.AddPassword(filePath, password, timestamp)
}

func GetPasswordList() []string {
	var displayList []string
	passwordList := data.GetPasswordList()
	for _, password := range passwordList {
		displayList = append(displayList, password.FilePath+" : "+password.Timestamp)
	}
	return displayList
}

func CreatePasswordEntry() *widget.Entry {
	return widget.NewEntry()
}

func CreatePasswordButton(filePath string, passwordEntry *widget.Entry) *widget.Button {
	return widget.NewButton("Добавить пароль", func() {
		password := passwordEntry.Text
		AddPasswordForFile(filePath, password)
	})
}
