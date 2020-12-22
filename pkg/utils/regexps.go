package utils

import "regexp"

var (
	usernamePattern    = regexp.MustCompile("^(?i)[a-z0-9_-]+$")
	emailPattern       = regexp.MustCompile("^(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])$")
	hexRGBColorPattern = regexp.MustCompile("^(?i)([0-9a-f]{6}|[0-9a-f]{3})$")
)

// IsValidUsername reports whether username is valid
func IsValidUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

// IsValidEmail reports whether email is valid
func IsValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

// IsValidHexRGBColor reports whether color is valid
func IsValidHexRGBColor(color string) bool {
	return hexRGBColorPattern.MatchString(color)
}
