package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secret = []byte("&@HKD^&*@#$")

type Claims struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Authority int    `json:"authority"`
	jwt.RegisteredClaims
}

// GenerateToken 签发token
func GenerateToken(id uint, username string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := jwt.NewNumericDate(nowTime.Add(24 * 30 * time.Hour))
	var claims = &Claims{
		Id:        id,
		Username:  username,
		Authority: authority,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expireTime,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

// ValidToken 验证用户Token
func ValidToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if tokenClaims == nil || err != nil || !tokenClaims.Valid {
		return nil, err
	}
	return tokenClaims.Claims.(*Claims), nil
}
