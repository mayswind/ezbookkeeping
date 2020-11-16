package utils

import "regexp"

var (
	UsernamePattern = regexp.MustCompile("^(?i)[a-z0-9_-]+$")
	EmailPattern = regexp.MustCompile("^(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])$")
	RGBColorPattern = regexp.MustCompile("^(?i)([0-9a-f]{6}|[0-9a-f]{3})$")
)

func IsValidUsername(username string) bool {
	return UsernamePattern.MatchString(username)
}

func IsValidEmail(email string) bool {
	return EmailPattern.MatchString(email)
}

func IsValidRGBColor(color string) bool {
	return RGBColorPattern.MatchString(color)
}
