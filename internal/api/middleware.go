package api

import (
	"DevIntApp/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

const prefix = "Bearer"

func (a *Application) RoleMiddleware(users ...ds.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := a.extractTokenFromHandler(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized1"})
			return

		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized2"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		log.Println(claims)

		userIDnew, ok := claims["userID"].(float64)
		log.Println(userIDnew, ok)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized4"})
			return
		}
		c.Set("userID", float64(userIDnew))

		user_role, ok := claims["isModerator"].(bool)
		log.Println(user_role, ok)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized5"})
			return
		}

		if !a.isRoleAllowed(user_role, users) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized6"})
			return
		}
		c.Next()

	}
}

func (a *Application) isRoleAllowed(userRole bool, users []ds.Users) bool {
	for _, v := range users {
		log.Println(v, userRole)
		if v.IsModerator == userRole {
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
