package api

import (
	"DevIntApp/internal/app/config"
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/dsn"
	"DevIntApp/internal/app/repository"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

func (a *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	var err error

	r.SetFuncMap(template.FuncMap{
		"replaceNewline": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "/n", "<br>"))
		},
	})

	r.Static("/css", "./resources")
	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		childmealsquery := c.Query("childmeal") // Получаем поисковый запрос из URL
		var FilteredMeals []ds.Meals

		if childmealsquery == "" {
			FilteredMeals, err = a.repo.GetAllMeals()
			if err != nil {
				log.Println("unable to get all meals")
				c.Error(err)
				return
			}
		} else {
			FilteredMeals, err = a.repo.GetMealByMealInfo(childmealsquery)
			if err != nil {
				log.Println("unable to get meal by info")
				FilteredMeals = []ds.Meals{}
			}
		}

		var milkreq_len int
		var milkreq_ID int
		milkreq_wrk, err := a.repo.GetWorkingMilkRequest()
		if err != nil {
			log.Println("unable to get working milk request")
		}
		if len(milkreq_wrk) == 0 {
			milkreq_len = 0
			milkreq_ID = 0

		} else {
			milkmeals_in_wrk_req, err := a.repo.GetMealsIDsByMilkRequestID(milkreq_wrk[0].ID)
			if err != nil {
				log.Println("unable to get meals ids by cart")
			}
			milkreq_len = len(milkmeals_in_wrk_req)
			milkreq_ID = milkreq_wrk[0].ID
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"filteredCards": FilteredMeals,
			"searchQuery":   childmealsquery,
			"meals_cnt":     milkreq_len,
			"milkreq_ID":    milkreq_ID,
		})
	})

	r.POST("/home", func(c *gin.Context) {

		id := c.PostForm("add")
		milkmeal_ID, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant transform ind", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		milkreq_wrk, err := a.repo.GetWorkingMilkRequest()
		var milkreq_ID int
		if len(milkreq_wrk) == 0 {
			new_milkreq, err := a.repo.CreateMilkRequest()
			if err != nil {
				log.Println("unable to create milk request")
			}
			milkreq_ID = new_milkreq[0].ID
		} else {
			milkreq_ID = milkreq_wrk[0].ID
		}

		a.repo.AddToMilkRequest(milkreq_ID, milkmeal_ID)

		c.Redirect(301, "/home")

	})

	r.GET("/meal/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		childmeal, err := a.repo.GetMilkRequestByID(index)
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "meal.html", gin.H{
			"title":     "Main website",
			"meal_data": childmeal,
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

		milk_req_status, err := a.repo.GetMilkRequestStatusByID(index)
		if err != nil {
			log.Printf("cant get milkreq by id %v", err)
		}
		if milk_req_status == 3 {
			c.Redirect(301, "/home")
		}

		MealsIDs, err := a.repo.GetMealsIDsByMilkRequestID(index)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")
			c.Error(err)
			return
		}

		MealsInMilkReq := []ds.Meals{}
		for _, v := range MealsIDs {
			meal_tmp, err := a.repo.GetMealByID(v.MealID)
			if err != nil {
				c.Error(err)
				return
			}
			MealsInMilkReq = append(MealsInMilkReq, meal_tmp[0])
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
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		a.repo.DeleteMilkRequest(index)
		c.Redirect(301, "/home")

	})

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
