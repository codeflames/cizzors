package utils

func PasswordLengthChecker(password string) bool {
	return len(password) >= 6
}

func UsernameLengthChecker(username string) bool {
	return len(username) >= 3
}
