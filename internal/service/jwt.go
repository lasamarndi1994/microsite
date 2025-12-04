package service

import (
	"micro-site/api/model"
	"micro-site/config"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user model.User) (string, error) {
	cfg := config.LoadConfig()
	secret := cfg.JWTSecretKey

	claims := jwt.MapClaims{
		"user_id":       user.Id,
		"uuid":          user.Uuid,
		"email":         user.Email,
		"mobile_number": user.MobileNumber,
		"exp":           time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
