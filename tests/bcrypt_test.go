package tests

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestBcrypt(t *testing.T) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("y18571139145.."), 12)

	fmt.Println(string(bytes))
}
