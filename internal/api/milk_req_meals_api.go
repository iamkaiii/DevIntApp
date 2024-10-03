package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) DeleteMealFromMilkReq(c *gin.Context) {
	var request schemas.DeleteMealFromMilkReqRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.DeleteMealFromMilkRequest(request.ID, request.MealID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was deleted from milk request")
}

func (a *Application) UpdateAmountMilkReqMeal(c *gin.Context) {
	var request schemas.UpdateAmountMilkReqMealRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateAmountMilkReqMeal(request.ID, request.MealID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal amount was changed in milk request")
}
