package api

import (
	"DevIntApp/internal/app/config"
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/dsn"
	"DevIntApp/internal/app/repository"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
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

	// API MILK REQUESTS

	r.GET("/api/meals", a.GetAllMeals)
	r.GET("/api/meal/:ID", a.GetMeal)
	r.POST("/api/meal", a.CreateMeal)                                                                                                        // да
	r.DELETE("/api/meal/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.DeleteMeal)                                                   // да
	r.PUT("/api/meal/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateMeal)                                                      // да
	r.POST("/api/meal_to_milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.AddMealToMilkReq) // да
	r.POST("/api/meal/change_pic/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.ChangePic)             // да

	r.GET("/api/milk_requests", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetAllMilkRequestsWithParams)
	r.GET("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetMilkRequest)
	r.PUT("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.UpdateFieldsMilkReq)
	r.DELETE("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.DeleteMilkRequest)
	r.PUT("/api/milk_request/form/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}, ds.Users{IsModerator: false}), a.FormMilkRequest)
	r.PUT("/api/milk_request/finish/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.FinishMilkRequest)

	r.DELETE("/api/milk_req_meals/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}, ds.Users{IsModerator: false}), a.DeleteMealFromMilkReq)
	r.PUT("/api/milk_req_meals/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateAmountMilkReqMeal)

	r.POST("/api/register_user", a.RegisterUser)
	r.POST("/api/login_user", a.LoginUser)
	r.POST("/api/logout", a.LogoutUser)
	r.PUT("/api/change_user_info", a.RoleMiddleware(ds.Users{IsModerator: true}, ds.Users{IsModerator: false}), a.ChangeUserInfo)

	r.GET("/protected", a.RoleMiddleware(ds.Users{IsModerator: true}), func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)
		c.JSON(http.StatusOK, gin.H{"message": "Пользователь авторизован с правами модератора", "userID": userID})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api/secret/get_meal_by_meal_info/:meal", a.GetMealsByName)
	var err error

	r.SetFuncMap(template.FuncMap{
		"replaceNewline": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "/n", "<br>"))
		},
		"replaceNewlineN": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "\n", "<br>"))
		},
	})

	r.Static("/css", "./resources")
	r.LoadHTMLGlob("templates/*")

	errRun := r.Run(":8001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if errRun != nil {
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

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(app.config.Minio.MinioAccess, app.config.Minio.MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *Application) UploadImage(c *gin.Context, image *multipart.FileHeader) (string, error) {
	openFile, err := image.Open()
	defer func() {
		err = openFile.Close()
		if err != nil {
			log.Println(err, "file")
		}
	}()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(openFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}
	reader := bytes.NewReader(fileBytes)
	_, err = a.minioClient.PutObject(context.Background(), a.config.Minio.BucketName, image.Filename, reader, image.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}
	url, err := a.minioClient.PresignedGetObject(context.Background(), a.config.Minio.BucketName, image.Filename, time.Second*24*60*60, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}
	return strings.Split(url.String(), "?")[0], nil
}

func (a *Application) DeleteImage(meal ds.Meals) error {
	splitedUrl := strings.Split(meal.ImageUrl, "/")
	err := a.minioClient.RemoveObject(context.Background(), a.config.Minio.BucketName, splitedUrl[len(splitedUrl)-1], minio.RemoveObjectOptions{})
	return err
}
