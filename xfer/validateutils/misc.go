package validateutils

import (
	"fmt"
	"github.com/ttacon/libphonenumber"
	"strings"
)

func ValidPhoneNumber(phone, country string) (bool, string) {
	if len(strings.TrimSpace(phone)) == 0 {
		return false, "Empty phone number"
	}
	_, err := libphonenumber.Parse(phone, country)
	if err != nil {
		if len(country) == 0 {
			return false, "Invalid phone number"
		} else {
			return false, fmt.Sprintf("Invalid phone number %s for country %s", phone, country)
		}
	}

	return true, ""
}
