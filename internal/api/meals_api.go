package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (a *Application) GetAllMeals(c *gin.Context) {
	var request schemas.GetAllMealsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meals, err := a.repo.GetAllMeals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	meals_cnt := len(meals)
	wrk_milk_req, err := a.repo.GetWorkingMilkRequest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var id int
	if len(wrk_milk_req) == 0 {
		id = 0
	} else {
		id = wrk_milk_req[0].ID
	}
	response := schemas.GetAllMealsResponse{ID: id, Count: meals_cnt, Meals: meals}
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetMeal(c *gin.Context) {
	var request schemas.GetMealRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal, err := a.repo.GetMealByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.GetMealResponse{Meal: meal}
	c.JSON(http.StatusOK, response)
}

func (a *Application) CreateMeal(c *gin.Context) {
	var request schemas.CreateMealRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.CreateMeal(request.Meals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Meal was created")
}

func (a *Application) DeleteMeal(c *gin.Context) {
	var request schemas.GetMealRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.DeleteMealByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was deleted")
}

func (a *Application) UpdateMeal(c *gin.Context) {
	var request schemas.UpdateMealRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request.Meals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request.Meals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateMealByID(request.ID, request.Meals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was updated")
}

func (a *Application) AddMealToMilkReq(c *gin.Context) {
	var request schemas.AddMealToMilkReqRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	new_milk_request, err := a.repo.CreateMilkRequest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	new_milk_request_id := new_milk_request.ID
	meal_id, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = a.repo.AddToMilkRequest(new_milk_request_id, meal_id)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was added")
}

func (a *Application) ChangePic(c *gin.Context) {
	var request schemas.ChangePicRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(request.ID, request.ImageUrl)
	err := a.repo.ChangePicByID(request.ID, request.ImageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal Pic was updated")
}
