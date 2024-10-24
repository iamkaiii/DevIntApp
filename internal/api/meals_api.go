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
// @Description Returns a list of all meals.
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
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
// @Security BearerAuth
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
// @Security BearerAuth
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
// @Description Delete meal using it's ID
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
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
	meal, err := a.repo.GetMealByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(meal)

	err = a.repo.DeleteMealByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	intID, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	meal, err1 := a.repo.GetMealByID(request.ID)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = a.DeleteImage(c, meal)
	if err != nil {
		response := schemas.DeleteMealResponse{ID: intID, MessageResponse: "Meal was deleted successfully"}
		c.JSON(418, response)
		return
	}
	err = a.repo.ChangePicByID(request.ID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.DeleteMealResponse{ID: intID, MessageResponse: "Meal was deleted successfully"}
	c.JSON(http.StatusOK, response)
}

// @Summary Update meal by ID
// @Description Update meal using it's ID with parametres
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
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

// @Summary Add meal to milk request
// @Description This endpoint allows you to add a meal to a milk request by it's ID.
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Meal ID"
// @Param request query schemas.AddMealToMilkReqRequest true "AddMealToMilkReqRequest"
// @Success 200 {object} schemas.AddMealToMilkReqResponse "Meal added successfully"
// @Failure 400 {object} schemas.ResponseMessage "Bad Request"
// @Failure 500 {object} schemas.ResponseMessage "Internal Server Error"
// @Router /api/meal_to_milk_request/{ID} [post]
func (a *Application) AddMealToMilkReq(c *gin.Context) {
	var request schemas.AddMealToMilkReqRequest
	idFromQuery := c.Param("ID")
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
	mealID, err := strconv.Atoi(idFromQuery)
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
	response := schemas.AddMealToMilkReqResponse{MealID: mealID, MilkRequestID: newMilkRequestID, MessageResponse: "Meal was added successfully to a new request"}
	c.JSON(http.StatusOK, response)
}

// @Summary Change picture By ID
// @Description Delete meal using it's ID
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Meal ID"
// @Param image formData file true "File"
// @Success 200 {object} schemas.ResponseMessage "Picture was changed sucessfully"
// @Router /api/meal/change_pic/{ID} [post]
func (a *Application) ChangePic(c *gin.Context) {
	var request schemas.ChangeImgRequest
	var err error
	request.ID = c.Param("ID")
	request.Image, err = c.FormFile("image")
	imageUrl, err := a.UploadImage(c, request.Image)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = a.repo.ChangePicByID(request.ID, imageUrl)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Meal Pic was updated")
}
