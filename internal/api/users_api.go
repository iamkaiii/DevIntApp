package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) RegisterUser(c *gin.Context) {
	var request schemas.RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Users.Login == "" || request.Users.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Login or Password can not be empty"})
		return
	}
	if len(request.Password) < 5 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password length must be more than 4 symbols"})
		return
	}

	err, res := a.repo.RegisterUser(request.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User with this login already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully!"})

}

func (a *Application) LoginUser(c *gin.Context) {
	var request schemas.LoginUserRequest
}
