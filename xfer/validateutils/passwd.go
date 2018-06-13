package validateutils

import "fmt"

var (
	MIN_PASSWORD_LEN = 10
	MAX_PASSWORD_LEN = 60
)

func ValidPassword(password string) (bool, string) {
	if len(password) < MIN_PASSWORD_LEN {
		return false, fmt.Sprintf("Password must be at least %d characters long", MIN_PASSWORD_LEN)
	} else if len(password) > MAX_PASSWORD_LEN {
		return false, fmt.Sprintf("Password must be less than %d characters long", MAX_PASSWORD_LEN)
	}
	return true, ""
}
