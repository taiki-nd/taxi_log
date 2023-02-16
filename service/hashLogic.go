package service

import "golang.org/x/crypto/bcrypt"

/**
 *
 */
func HashMethod(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), 8)
	if err != nil {
		return "", err
	}
	return string(hashedValue), err
}
