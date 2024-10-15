package api

import (
	"DevIntApp/internal/app/config"
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/dsn"
	"DevIntApp/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

// @title DevIntApp
// @version 1.1
// @description This is API for Milk Kitchen requests
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func (a *Application) Run() {
	log.Println("Server start up")
	r := gin.Default()

	r.GET("/api/meals", a.GetAllMeals)                          // да
	r.GET("/api/meal/:ID", a.GetMeal)                           // да + картинка?
	r.POST("/api/meal", a.CreateMeal)                           // да
	r.DELETE("/api/meal/:ID", a.DeleteMeal)                     // да
	r.PUT("/api/meal/:ID", a.UpdateMeal)                        // да
	r.POST("/api/meal_to_milk_request/:ID", a.AddMealToMilkReq) // да
	r.POST("/api/meal/change_pic/:ID", a.ChangePic)             // да

	r.GET("/api/milk_requests", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetAllMilkRequestsWithParams)
	r.GET("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetMilkRequest)
	r.PUT("/api/milk_request/:ID", a.UpdateFieldsMilkReq)
	r.DELETE("/api/milk_request/:ID", a.DeleteMilkRequest)
	r.PUT("/api/milk_request/form/:ID", a.FormMilkRequest)
	r.PUT("/api/milk_request/finish/:ID", a.FinishMilkRequest)

	r.DELETE("/api/milk_req_meals/:ID", a.DeleteMealFromMilkReq)
	r.PUT("/api/milk_req_meals/:ID", a.UpdateAmountMilkReqMeal)

	///ЛАБА 4 ПО РИПУ///

	r.POST("/api/register_user", a.RegisterUser)
	r.POST("/api/login_user", a.LoginUser)

	r.GET("/protected", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)
		c.JSON(http.StatusOK, gin.H{"message": "Пользователь авторизован с правами модератора", "userID": userID})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/css", "./resources")
	err := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
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
