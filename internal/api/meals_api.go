package api

import (
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Get all meals
// @Description Returns a list of all meals.
// @Tags meals
// @Accept json
// @Produce json
// @Success 200 {object} schemas.ResponseMessage "List of meals retrieved successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/meals [get]
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
	mealsCnt := len(meals)
	activeMilkRequest, err := a.repo.GetWorkingMilkRequest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var id int
	if len(activeMilkRequest) == 0 {
		c.JSON(http.StatusOK, gin.H{"Count": mealsCnt, "Meals": meals})
		return
	} else {
		id = activeMilkRequest[0].ID
	}
	response := schemas.GetAllMealsResponse{ID: id, Count: mealsCnt, Meals: meals}
	c.JSON(http.StatusOK, response)
	return
}

// GetMeal godoc
// @Summary Get a meal by ID
// @Description Get details of a meal using its ID
// @Tags meals
// @Accept json
// @Produce json
// @Param ID path string true "Meal ID"
// @Success 200 {object} schemas.GetMealResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/meal/{ID} [get]
func (a *Application) GetMeal(c *gin.Context) {
	var request schemas.GetMealRequest
	request.ID = c.Param("ID")
	log.Println(request.ID)
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
	log.Println(request)
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
	newMilkRequest, err := a.repo.CreateMilkRequest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newMilkRequestID := newMilkRequest.ID
	mealID, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = a.repo.AddToMilkRequest(newMilkRequestID, mealID)
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
