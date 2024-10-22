package api

import (
	"DevIntApp/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const prefix = "Bearer"

func (a *Application) RoleMiddleware(allowedRoles ...ds.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		as := c.GetHeader("Authorization")
		log.Println(as)
		tokenString := a.extractTokenFromHandler(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}
		log.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		userID, ok := claims["userID"].(float64)

		if a.tokenBlacklist(userID, tokenString) {
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

func (a *Application) tokenBlacklist(userID float64, token string) bool {
	userIDStr := strconv.FormatFloat(userID, 'f', 0, 64)
	res, err := a.repo.CheckBlacklist(userIDStr)
	if err != nil {
		log.Println(err)
		return false
	}
	if res == "revoked" {
		return true
	}
	return false
}
