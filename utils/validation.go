package utils

import "unicode"

// ValidateItemName checks that all characters are letters or spaces
func ValidateItemName(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' ' {
			return false
		}
	}
	return true
}
