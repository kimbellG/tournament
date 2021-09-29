package controller

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-=+_{}[];,."
const (
	letterIDxBits = 6
	letterIDxMask = 1<<letterIDxBits - 1
	letterIDxMax  = 63 / letterIDxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for i, cache, remain := n-1, src.Int63(), letterIDxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDxMax
		}

		if idx := int(cache & letterIDxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}

		cache >>= letterIDxBits
		remain--
	}

	return sb.String()
}

func generatePassword() string {
	const (
		minSize = 8
		maxSize = 20
	)

	password := generateRandomString(minSize + rand.Intn(maxSize-minSize))

	return password
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
