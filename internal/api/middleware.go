package api

import (
	"DevIntApp/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const prefix = "Bearer"

func (a *Application) RoleMiddleware(allowedRoles ...ds.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := a.extractTokenFromHandler(c.Request)
		log.Println(tokenString)
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
		userID, ok := claims["userID"].(float64)

		if !a.tokenActive(userID, tokenString) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен устарел"})
			return
		}

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
		c.Set("isModerator", userRole)

		if !isRoleAllowed(allowedRoles, userRole) { //проверка, является ли модератором
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отказано в доступе"})
			return
		}
		c.Next()
	}
}

func isRoleAllowed(roles []ds.Users, userRole bool) bool {
	for _, allowedRole := range roles {
		if userRole == allowedRole.IsModerator {
			return true
		}
	}
	return false
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

func (a *Application) tokenActive(userID float64, token string) bool {
	userIDStr := strconv.FormatFloat(userID, 'f', 0, 64)
	res, err := a.repo.CheckActive(userIDStr)
	if err != nil {
		return false
	}
	if res == "revoked" || res != token {
		return false
	}
	return true
}
