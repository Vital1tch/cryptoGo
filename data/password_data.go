package data

type PasswordInfo struct {
	FilePath  string
	Password  string
	Timestamp string
}

var passwordList []PasswordInfo

func AddPassword(filePath, password, timestamp string) {
	passwordList = append(passwordList, PasswordInfo{
		FilePath:  filePath,
		Password:  password,
		Timestamp: timestamp,
	})
}

func GetPasswordList() []PasswordInfo {
	return passwordList
}
