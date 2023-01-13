package utils

import (
	"crypto/rand"
	"github.com/syyongx/php2go"
)

func (u *Utils) Bin2hex(str string) (string, error) {
	return php2go.Bin2hex(str)
}

func (u *Utils) RandomBytes(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, randomBytesErr := rand.Read(randomBytes)
	if randomBytesErr != nil {
		return "", randomBytesErr
	}

	return string(randomBytes), nil
}
