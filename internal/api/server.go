package api

import (
	"DevIntApp/internal/app/config"
	"DevIntApp/internal/app/dsn"
	"DevIntApp/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func (a *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/api/meals", a.GetAllMeals)
	r.GET("/api/meal/:ID", a.GetMeal)
	r.POST("/api/meal", a.CreateMeal)
	r.DELETE("/api/meal/:ID", a.DeleteMeal)
	r.PUT("/api/meal/:ID", a.UpdateMeal)
	r.POST("/api/meal_to_milkreq/:ID", a.AddMealToMilkReq)
	r.POST("api/meal/change_pic/:ID", a.ChangePic)

	r.GET("/api/milk_requests", a.GetAllMilkRequestsWithParams)
	r.GET("/api/milk_request/:ID", a.GetMilkRequest)
	r.PUT("/api/milk_request/:ID", a.UpdateFieldsMilkReq)
	r.DELETE("/api/milk_request/:ID", a.DeleteMilkRequest)
	r.PUT("/api/milk_request/form/:ID", a.FormMilkRequest)
	r.PUT("/api/milk_request/finish/:ID", a.FinishMilkRequest)

	r.DELETE("/api/milk_req_meals/:ID", a.DeleteMealFromMilkReq)
	r.PUT("/api/milk_req_meals/:ID", a.UpdateAmountMilkReqMeal)

	r.POST("/api/user_reg", a.CreateUser)

	r.Static("/css", "./resources")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
