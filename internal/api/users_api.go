package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) CreateUser(c *gin.Context) {
	var request schemas.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.CreateUser(request.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "User was created")
}
