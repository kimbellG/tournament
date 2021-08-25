package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type passwordBuffer []string

func TestPasswordGenerate(t *testing.T) {
	buf := passwordBuffer{}

	for i := 8; i < 20; i++ {
		password := generateRandomString(i)
		t.Logf("password from func: %v", password)

		assert.Equal(t, i, len(password), "length of password should be n")
		assert.Falsef(t, buf.isIn(password), "password(%v) shouldn't be in buf(%v)", password, buf)
	}

}

func (pb passwordBuffer) isIn(password string) bool {
	for _, pas := range pb {
		if pas == password {
			return true
		}
	}

	return false
}
