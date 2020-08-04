package handler

import (
	"errors"
	"regexp"
)

const (
	EMAIL_REGEXP = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

func validateEmail(email string) error  {
	regexp := regexp.MustCompile(EMAIL_REGEXP)
	if !regexp.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}