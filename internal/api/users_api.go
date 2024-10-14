package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// @Summary Register a new user
// @Description Register a new user with provided JSON data
// @Tags users
// @Accept json
// @Produce json
// @Param body schemas.RegisterUserRequest true "User registration data"
// @Success 200 {object} schemas.ResponseMessage "User created successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"

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

// @Summary Login a user
// @Description Authenticates a user and returns a JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param body schemas.LoginUserRequest true "User login data"
// @Success 200 {object} schemas.TokenResponse "User logged in successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"

func (a *Application) LoginUser(c *gin.Context) {
	var request schemas.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Users.Login == "" || request.Users.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Login or Password can not be empty"})
		return
	}
	err, token := a.repo.LoginUser(request.Users)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
