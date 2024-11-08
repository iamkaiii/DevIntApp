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
	"strconv"
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
	r.POST("/api/meal", a.RoleMiddleware(ds.Users{IsModerator: true}), a.CreateMeal)                                                         // да
	r.DELETE("/api/meal/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.DeleteMeal)                                                   // да
	r.PUT("/api/meal/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateMeal)                                                      // да
	r.POST("/api/meal_to_milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.AddMealToMilkReq) // да
	r.POST("/api/meal/change_pic/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.ChangePic)                                           // да

	r.GET("/api/milk_requests", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetAllMilkRequestsWithParams)
	r.GET("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetMilkRequest)
	r.PUT("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.UpdateFieldsMilkReq)
	r.DELETE("/api/milk_request/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.DeleteMilkRequest)
	r.PUT("/api/milk_request/form/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.FormMilkRequest)
	r.PUT("/api/milk_request/finish/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.FinishMilkRequest)

	r.DELETE("/api/milk_req_meals/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.DeleteMealFromMilkReq)
	r.PUT("/api/milk_req_meals/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateAmountMilkReqMeal)

	r.POST("/api/register_user", a.RegisterUser)
	r.POST("/api/login_user", a.LoginUser)
	r.POST("/api/logout", a.LogoutUser)

	r.GET("/protected", a.RoleMiddleware(ds.Users{IsModerator: true}), func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)
		c.JSON(http.StatusOK, gin.H{"message": "Пользователь авторизован с правами модератора", "userID": userID})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	r.GET("/home", func(c *gin.Context) {
		childMealsQuery := c.Query("childmeal") // Получаем поисковый запрос из URL
		var FilteredMeals []ds.Meals

		if childMealsQuery == "" {
			FilteredMeals, err = a.repo.GetAllMeals()
			if err != nil {
				log.Println("unable to get all meals")
				return
			}
		} else {
			FilteredMeals, err = a.repo.GetMealsByMealInfo(childMealsQuery)
			if err != nil {
				log.Println("unable to get meal by info")
				FilteredMeals = []ds.Meals{}
			}
		}

		var milkReqLen int
		var milkReqID int
		milkReqWorking, err := a.repo.GetWorkingMilkRequest()
		if err != nil {
			log.Println("unable to get working milk request")
		}
		if len(milkReqWorking) == 0 {
			milkReqLen = 0
			milkReqID = 0

		} else {
			MilkMealsInWorkingReq, err1 := a.repo.GetMealsIDsByMilkRequestID(milkReqWorking[0].ID)
			if err1 != nil {
				log.Println("unable to get meals ids by cart")
			}
			milkReqLen = len(MilkMealsInWorkingReq)
			milkReqID = milkReqWorking[0].ID

		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"filteredCards": FilteredMeals,
			"searchQuery":   childMealsQuery,
			"meals_cnt":     milkReqLen,
			"milkreq_ID":    milkReqID,
		})
	})

	r.POST("/home", func(c *gin.Context) {

		id := c.PostForm("add")
		milkMealID, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Println("cant transform ind", err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		milkReqWorking, err := a.repo.GetWorkingMilkRequest()
		var milkReqID int
		if len(milkReqWorking) == 0 {
			newMilkReq, err := a.repo.CreateMilkRequest()
			if err != nil {
				log.Println("unable to create milk request")
			}
			milkReqID = newMilkReq.ID
		} else {
			milkReqID = milkReqWorking[0].ID
		}

		err = a.repo.AddToMilkRequest(milkReqID, milkMealID)

		c.Redirect(301, "/home")

	})

	r.GET("/meal/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL

		childMeal, err := a.repo.GetMealByID(id)
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "meal.html", gin.H{
			"title":     "Main website",
			"meal_data": childMeal,
		})
	})

	r.GET("/milkreq/:id", func(c *gin.Context) {

		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get milkreq by id %v", err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		milkReqStatus, err := a.repo.GetMilkRequestStatusByID(index)
		if err != nil {
			log.Printf("cant get milkreq by id %v", err)
		}
		if milkReqStatus == 3 {
			c.Redirect(301, "/home")
		}

		MealsIDs, err := a.repo.GetMealsIDsByMilkRequestID(index)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")

			return
		}

		MealsInMilkReq := []ds.Meals{}
		for _, v := range MealsIDs {
			vString := strconv.Itoa(v)
			mealTemp, err := a.repo.GetMealByID(vString)
			if err != nil {
				return
			}
			MealsInMilkReq = append(MealsInMilkReq, mealTemp)
		}

		c.HTML(http.StatusOK, "milkreq.html", gin.H{
			"title":          "Корзина",
			"MealsInMilkReq": MealsInMilkReq,
			"MilkReqID":      index,
		})
	})

	r.POST("/milkreq/:id", func(c *gin.Context) {
		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get cart by id %v", err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		err = a.repo.DeleteMilkRequest(index)
		if err != nil {
			log.Println("unable to delete milk request")
			return
		}
		c.Redirect(301, "/home")

	})

	r.GET("/begin", func(c *gin.Context) {
		description := "Услуга позволяет самостоятельно \n(минуя кабинет врача):\n- заказывать питание на молочной кухне;\n- изменять пункт выдачи продуктов;\n- управлять графиком получения продуктов питания;\n- просматривать информацию о полученной продукции."
		c.HTML(http.StatusOK, "begin.html", gin.H{
			"PageTitle":   "Заказ питания",
			"Description": description,
		})
	})

	errRun := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
	log.Println("uploaded")
	url, err := a.minioClient.PresignedGetObject(context.Background(), a.config.Minio.BucketName, image.Filename, time.Second*24*60*60, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}
	return strings.Split(url.String(), "?")[0], nil
}

func (a *Application) DeleteImage(meal ds.Meals) error {
	splitedUrl := strings.Split(meal.ImageUrl, "/")
	log.Println(splitedUrl)
	err := a.minioClient.RemoveObject(context.Background(), a.config.Minio.BucketName, splitedUrl[len(splitedUrl)-1], minio.RemoveObjectOptions{})
	return err
}
