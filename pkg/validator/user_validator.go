package validator

import "regexp"

const (
	emailPattern    = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	usernamePattern = "^[A-Za-z][A-Za-z0-9_]{7,29}$"
	passwordPattern = "^[A-Za-z][A-Za-z0-9_]{7,29}$"
)

func ValidEmail(email string) bool {
	matched, _ := regexp.Match(emailPattern, []byte(email))
	return matched
}

func ValidUsername(username string) bool {
	matched, _ := regexp.Match(usernamePattern, []byte(username))
	return matched
}

func ValidPassword(password string) bool {
	matched, _ := regexp.Match(passwordPattern, []byte(password))
	return matched
}
