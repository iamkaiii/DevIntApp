package api

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (a *Application) GetAllMilkRequestsWithParams(c *gin.Context) {
	var request schemas.GetAllMilkRequestsWithParamsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.FromDate.IsZero() {
		request.FromDate = time.Date(2000, time.January, 1, 0, 0, 0, 396641000, time.UTC)
	}
	if request.ToDate.IsZero() {
		request.ToDate = time.Now()
	}
	milk_requests, err := a.repo.GetAllMilkRequestsWithFilters(request.Status, request.HavingStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.GetAllMilkRequestsWithParamsResponse{MilkRequests: milk_requests}
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetMilkRequest(c *gin.Context) {
	var request schemas.GetMilkRequestRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	milk_request, err := a.repo.GetMilkRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	meals_ids_in_request, err := a.repo.GetMealsIDsByMilkRequestID(milk_request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	meals := make([]ds.Meals, 0, len(meals_ids_in_request))
	for _, v := range meals_ids_in_request {
		v_string := strconv.Itoa(v)
		meal_to_append, err := a.repo.GetMealByID(v_string)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		meals = append(meals, meal_to_append)
	}
	response := schemas.GetMilkRequestResponse{MilkRequest: milk_request, MilkRequestMeals: meals}
	c.JSON(http.StatusOK, response)
}
