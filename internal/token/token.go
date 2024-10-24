package token

import (
	"DevIntApp/internal/app/ds"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func GenerateJWTToken(user ds.Users) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["isModerator"] = user.IsModerator
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
