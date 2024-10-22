package api

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// @Summary Register a new user
// @Description Registers a new user.
// @Tags users
// @Accept json
// @Produce json
// @Param  body   body   schemas.RegisterUserRequest true "User registration data"
// @Success 200 {object} schemas.ResponseMessage "User registered successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 409 {object} schemas.ResponseMessage "User with this login already exists"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router      /api/register_user [post]
func (a *Application) RegisterUser(c *gin.Context) {
	var request schemas.RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Login == "" || request.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Login or Password can not be empty"})
		return
	}
	if len(request.Password) < 5 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password length must be more than 4 symbols"})
		return
	}
	requestTemp := ds.Users{Login: request.Login, Password: request.Password}
	err, res := a.repo.RegisterUser(requestTemp)
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
// @Param body body schemas.LoginUserRequest true "User login data"
// @Success 200 {object} schemas.ResponseMessage "User logged in successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router      /api/login_user [post]
func (a *Application) LoginUser(c *gin.Context) {
	var request schemas.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Login == "" || request.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Login or Password can not be empty"})
		return
	}
	requestTemp := ds.Users{Login: request.Login, Password: request.Password}
	err, token := a.repo.LoginUser(requestTemp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Logout
// @Description Log out the user by blacklisting the token
// @Tags users
// @Accept json
// @Produce json
// @Param body body schemas.LogoutUserRequest true "User login data"
// @Success 200 {object} schemas.ResponseMessage "User logged out successfully"
// @Failure 401 {object} schemas.ResponseMessage "Missing token"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/logout [post]
func (a *Application) LogoutUser(c *gin.Context) {
	var request schemas.LogoutUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(request.Login)
	err := a.repo.LogoutUser(request.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
