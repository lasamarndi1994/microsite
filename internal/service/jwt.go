package service

import (
	"micro-site/api/model"
	"micro-site/config"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
* Generate JWT token for user
* @param user model.User
* @return string, error
 */
func GenerateJWT(user model.User) (string, error) {
	cfg := config.LoadConfig()
	secret := cfg.JWTSecretKey

	claims := jwt.MapClaims{
		"user_id":       user.Id,
		"uuid":          user.Uuid,
		"email":         user.Email,
		"mobile_number": user.MobileNumber,
		"exp":           time.Now().Add(time.Hour * 2).Unix(), // expires in 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

/*
* Generate JWT token for user
* @param user model.User
* @return string, error
 */
func GenerateAdminJWT(admin model.Admin) (string, error) {
	cfg := config.LoadConfig()
	secret := cfg.JWTSecretKey

	claims := jwt.MapClaims{
		"user_id": admin.Id,
		"email":   admin.Email,
		"exp":     time.Now().Add(time.Hour * 2).Unix(), // expires in 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
