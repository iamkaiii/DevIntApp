package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

const prefix = "Bearer"

func (a *Application) RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := a.extractTokenFromHandler(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return

		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		log.Println(claims)

		userID, ok := claims["userID"].(float64)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}
		c.Set("userID", float64(userID))

		userRole, ok := claims["isModerator"].(bool)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		if userRole != true { //проверка, является ли модератором
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отказано в доступе"})
			return
		}
		c.Next()

	}
}

func (a *Application) extractTokenFromHandler(req *http.Request) string {
	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}
	return strings.Split(bearerToken, " ")[1]
}
