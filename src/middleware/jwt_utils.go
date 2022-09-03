package middleware

import (
	"fmt"
	"time"

	"github.com/fnmzgdt/e_shop/src/utils"
	"github.com/golang-jwt/jwt"
)

func NewJWT(ttl time.Duration, content interface{}) (string, error) {
	privateKey := utils.GetEnv("JWT_SALT", "")
	now := time.Now()

	claims := make(jwt.MapClaims)
	claims["data"] = content            // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	return token, nil
}

func Validate(token string) (interface{}, error) {
	//jwt.ParseWithClaims()
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		privateKey := utils.GetEnv("JWT_SALT", "")
		return []byte(privateKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return claims["data"], nil
}
