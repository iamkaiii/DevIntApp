package api

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// GetAllMilkRequestsWithParams godoc
// @Summary Получить все заявки на молочную кухню с параметрами
// @Description Получить список запросов на молоко с возможностью фильтрации по статусу и датам
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param status query string false "Статус заявки"
// @Param is_status query string false "Наличие статуса"
// @Success 200 {object} schemas.GetAllMilkRequestsWithParamsResponse
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/milk_requests [get]
// @Security BearerAuth
func (a *Application) GetAllMilkRequestsWithParams(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("status"))
	havingStatus, _ := strconv.ParseBool(c.Query("is_status"))

	isModerator, ok := c.Get("isModerator")
	if !ok {
		c.JSON(http.StatusInternalServerError, ok)
		return
	}
	isModeratorBool := isModerator.(bool)

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, ok)
		return
	}
	userIDInt := userID.(float64)

	milkRequests, err := a.repo.GetAllMilkRequestsWithFilters(status, havingStatus, isModeratorBool, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := schemas.GetAllMilkRequestsWithParamsResponse{MilkRequests: milkRequests}
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetMilkRequest(c *gin.Context) {
	var request schemas.GetMilkRequestRequest
	request.ID = c.Param("ID")
	IntID, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("error was there")
		return
	}
	MilkRequest, err := a.repo.GetMilkRequestByID(IntID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	MealsIDsInRequest, err := a.repo.GetMealsIDsByMilkRequestID(MilkRequest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	meals := make([]ds.Meals, 0, len(MealsIDsInRequest))
	for _, v := range MealsIDsInRequest {
		vString := strconv.Itoa(v)
		MealsToAppend, err := a.repo.GetMealByID(vString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		meals = append(meals, MealsToAppend)
	}
	response := schemas.GetMilkRequestResponse{MilkRequest: MilkRequest, Count: len(MealsIDsInRequest), MilkRequestMeals: meals}
	c.JSON(http.StatusOK, response)
}

func (a *Application) UpdateFieldsMilkReq(c *gin.Context) {
	var request schemas.UpdateFieldsMilkReqRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateFieldsMilkReq(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Fields was updated")
}

func (a *Application) DeleteMilkRequest(c *gin.Context) {
	var request schemas.DeleteMilkRequestRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = a.repo.DeleteMilkRequest(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "MilkRequest was deleted")
}

func (a *Application) FormMilkRequest(c *gin.Context) {
	var request schemas.FormMilkRequestRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.FormMilkRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Milk Request was Formed")
}

func (a *Application) FinishMilkRequest(c *gin.Context) {
	var request schemas.FinishMilkRequestRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.FinishMilkRequest(id, request.Status, request.DeliveryDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Milk Request was Finished")
}
