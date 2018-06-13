package validateutils

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	diacritics        = ".,:;!?'\""
	digits            = "1234567890"
	basicLatinLetters = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
)

func ContainsOnly(str, validChars string) (bool, string) {
	for _, r := range str {
		if !strings.ContainsRune(validChars, r) {
			return false, fmt.Sprintf(`Must contain only one of the following characters "%s"`, validChars)
		}
	}
	return true, ""
}

func ContainsOnlyDigits(str string) (bool, string) {
	ok, _ := ContainsOnly(str, digits)
	if !ok {
		return false, "Must contain only digits"
	}
	return true, ""
}

func ContainsOnlyLettersNumbersSpaces(str string) (bool, string) {
	return ContainsOnlyLettersNumbersSpacesOr(str, "")
}

func ContainsOnlyBasicLatinLettersAndSpaces(str string) (bool, string) {
	for _, r := range str {
		if !unicode.IsSpace(r) && !strings.ContainsRune(basicLatinLetters, r) {
			return false, "Must contain only basic latin letters, numbers or spaces"
		}
	}
	return true, ""
}

func ContainsOnlyLettersNumbersSpacesInterpunctions(str string) (bool, string) {
	return ContainsOnlyLettersNumbersSpacesOr(str, diacritics)
}

func ContainsOnlyLettersNumbersSpacesOr(str string, additionalChars string) (bool, string) {
	for _, r := range str {
		if !unicode.IsSpace(r) && !unicode.IsLetter(r) && !unicode.IsNumber(r) && !strings.ContainsRune(additionalChars, r) {
			if len(additionalChars) == 0 {
				return false, "Must contain only letters, numbers or spaces"
			} else {
				return false, `Must contain only letters, numbers, spaces or characters "` + additionalChars + `"`
			}
		}
	}
	return true, ""
}

func ContainsOnlyLettersNumbersOr(str string, additionalChars string) (bool, string) {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !strings.ContainsRune(additionalChars, r) {
			if len(additionalChars) == 0 {
				return false, "Must contain only letters, numbers or spaces"
			} else {
				return false, `Must contain only letters, numbers, spaces or characters "` + additionalChars + `"`
			}
		}
	}
	return true, ""
}

func ContainsOnlyLetterNumbersOrDiacritics(str string) (bool, string) {
	return ContainsOnlyLettersNumbersOr(str, diacritics)
}

func ValidSkype(str string) (bool, string) {
	return ContainsOnlyLettersNumbersOr(str, "._-")
}
