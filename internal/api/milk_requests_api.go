package api

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// GetAllMilkRequestsWithParams godoc
// @Summary Получить все заявки на молочную кухню с параметрами
// @Description Получить список запросов на молоко с возможностью фильтрации по статусу и пользователю
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param status query int true "Статус заявки" // Параметр query для фильтрации по статусу
// @Success 200 {object} schemas.GetAllMilkRequestsWithParamsResponse
// @Failure 400 {object} schemas.ResponseMessage "Неверный запрос, отсутствует параметр status или он некорректен"
// @Failure 500 {object} schemas.ResponseMessage "Ошибка сервера"
// @Router /api/milk_requests [get]
// @Security BearerAuth
func (a *Application) GetAllMilkRequestsWithParams(c *gin.Context) {
	// Извлекаем параметр "status" из query строки
	status := c.DefaultQuery("status", "") // Если параметр не найден, возвращаем пустую строку

	statusInt, err := strconv.Atoi(status) // Преобразуем status в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be an integer"})
		return
	}

	log.Println("Status:", statusInt)

	// Получаем userID из контекста
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}
	userIDInt := userID.(float64)

	// Получаем список заявок с фильтрацией по статусу и userID
	milkRequests, err := a.repo.GetAllMilkRequestsWithFilters(statusInt, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := schemas.GetAllMilkRequestsWithParamsResponse{MilkRequests: milkRequests}
	c.JSON(http.StatusOK, response)
}

// @Summary Получить данные заявки на молочную кухню
// @Description Получить данные заявки на молочную кухню по её ID, включая список связанных блюд
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param ID path int true "ID заявки"
// @Success 200 {object} schemas.GetMilkRequestResponse
// @Failure 400 {object} schemas.ResponseMessage "Неверный ID"
// @Failure 500 {object} schemas.ResponseMessage "Ошибка сервера"
// @Router /api/milk_request/{ID} [get]
// @Security BearerAuth
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
	log.Println(request)
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

// @Summary Удалить заявку на молочную кухню
// @Description Удаляет заявку на молочную кухню по её ID
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param ID path int true "ID заявки"  // ID заявки, передаваемое в пути
// @Param request query schemas.DeleteMilkRequestRequest false "Параметры запроса" // Параметры запроса (если есть)
// @Success 200 {string} string "MilkRequest was deleted"  // Сообщение об успешном удалении
// @Failure 400 {object} schemas.ResponseMessage "Неверный ID"  // Ошибка при неверном ID
// @Failure 500 {object} schemas.ResponseMessage "Ошибка сервера"  // Ошибка сервера при удалении заявки
// @Router /api/milk_request/{ID} [delete]  // URL-метод для удаления заявки
// @Security BearerAuth  // Требуется аутентификация
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

// @Summary Создать заявку на молочную кухню
// @Description Формирует заявку на молочную кухню по переданному ID и параметрам запроса
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param ID path int true "ID заявки"
// @Param queryParams query schemas.FormMilkRequestRequest true "Параметры заявки"
// @Success 200 {string} string "Milk Request was Formed"
// @Failure 400 {object} schemas.ResponseMessage "Ошибка в параметрах запроса"
// @Failure 500 {object} schemas.ResponseMessage "Ошибка сервера"
// @Router /api/milk_request/form/{ID} [put]
// @Security BearerAuth
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

// @Summary Завершить заявку на молочную кухню
// @Description Завершаем заявку на молочную кухню по переданному ID и параметрам в теле запроса (статус и дата доставки)
// @Tags milk_requests
// @Accept json
// @Produce json
// @Param ID path int true "ID заявки"
// @Param requestBody body schemas.FinishMilkRequestRequest true "Данные для завершения заявки"
// @Success 200 {string} string "Milk Request was Finished"
// @Failure 400 {object} schemas.ResponseMessage "Ошибка в параметрах запроса"
// @Failure 500 {object} schemas.ResponseMessage "Ошибка сервера"
// @Router /api/milk_request/finish/{ID} [put]
// @Security BearerAuth
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
	status := 1
	deliveryDate := time.Now()
	err := a.repo.FinishMilkRequest(id, status, deliveryDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Milk Request was Finished")
}
