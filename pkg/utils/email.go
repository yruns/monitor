package utils

import (
	"math/rand"
)

var template = "1234567890QWERTYUIOPASDFGHJKLZXCVBNM"

func GenerateEmailVerifyCode() string {
	templateLen := len(template)

	res := ""
	for i := 0; i < 4; i++ {
		res += string(template[rand.Intn(templateLen)])
	}

	return res
}
