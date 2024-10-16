package api

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Get all meals
// @Description Returns a list ofall meals.
// @Tags meals
// @Accept json
// @Produce json
// @Success 200 {object} schemas.GetAllMealsResponse "List of meals retrieved successfully"
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

// @Summary Get meal by ID
// @Description Get info about meal using its ID
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

// @Summary Create meal
// @Description Create meal with properties
// @Tags meals
// @Accept json
// @Produce json
// @Param body body schemas.CreateMealRequest true "Meal data"
// @Success 201 {object} schemas.CreateMealResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/meal [post]
func (a *Application) CreateMeal(c *gin.Context) {
	var request schemas.CreateMealRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := ds.Meals{MealInfo: request.MealInfo,
		MealWeight: request.MealWeight,
		MealBrand:  request.MealBrand,
		MealDetail: request.MealDetail,
		Status:     request.Status,
		ImageUrl:   request.ImageUrl,
	}
	ID, err := a.repo.CreateMeal(meal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.CreateMealResponse{
		ID:              ID,
		MessageResponse: "Meal was created successfully",
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary Delete meal by ID
// @Description Delete meal using its ID
// @Tags meals
// @Accept json
// @Produce json
// @Param ID path string true "Meal ID"
// @Success 200 {object} schemas.DeleteMealResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/meal/{ID} [delete]
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
	stringID, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.DeleteMealResponse{ID: stringID, MessageResponse: "Meal was deleted successfully"}
	c.JSON(http.StatusOK, response)
}

// @Summary Update meal by ID
// @Description Update meal using its ID with parametres
// @Tags meals
// @Accept json
// @Produce json
// @Param ID path string true "Meal ID"
// @Param body body schemas.UpdateMealRequest true "Update meal data"
// @Success 200 {object} schemas.UpdateMealResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/meal/{ID} [put]
func (a *Application) UpdateMeal(c *gin.Context) {
	var request schemas.UpdateMealRequest
	idFromRequest := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := ds.Meals{
		MealInfo:   request.MealInfo,
		MealWeight: request.MealWeight,
		MealBrand:  request.MealBrand,
		MealDetail: request.MealDetail,
		Status:     request.Status,
		ImageUrl:   request.ImageUrl,
	}
	err := a.repo.UpdateMealByID(idFromRequest, meal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	idString, err := strconv.Atoi(idFromRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.UpdateMealResponse{ID: idString, MessageResponse: "Meal was updated successfully"}
	c.JSON(http.StatusOK, response)
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
