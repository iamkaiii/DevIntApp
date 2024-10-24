package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Удалить блюдо из запроса на молоко
// @Description Удаляет блюдо из запроса на молоко по ID запроса и MealID
// @Tags meals_and_requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "ID заявки"
// @Param body body schemas.DeleteMealFromMilkReqRequest true "Delete meal from meal request"
// @Success 200 {string} string "Meal was deleted from milk request"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/milk_req_meals/{ID} [delete]
func (a *Application) DeleteMealFromMilkReq(c *gin.Context) {
	var request schemas.DeleteMealFromMilkReqRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.DeleteMealFromMilkRequest(id, request.MealID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was deleted from milk request")
}

// @Summary Обновить количество продуктов в заявке
// @Description Обновляет количество продукта в конкретной заявке
// @Tags meals_and_requests
// @Accept json
// @Produce json
// @Param ID path string true "Milk request ID"
// @Param body body schemas.UpdateAmountMilkReqMealRequest true "Update amount meals in milk request"
// @Success 200 {string} string "Meal was deleted from milk request"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/milk_req_meals/{ID} [put]
func (a *Application) UpdateAmountMilkReqMeal(c *gin.Context) {
	var request schemas.UpdateAmountMilkReqMealRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateAmountMilkReqMeal(id, request.MealID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal amount was changed in milk request")
}
