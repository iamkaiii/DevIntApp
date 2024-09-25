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
				c.Error(err)
				return
			}
		}

		var cart_id int

		wrk_cart, err := a.repo.GetWorkingCart()
		log.Println(err)
		if len(wrk_cart) == 0 {
			new_cart, err := a.repo.CreateCart()
			cart_id = new_cart[0].ID
			log.Println(err)
		} else {
			cart_id = wrk_cart[0].ID
		}

		MealsIDs, err := a.repo.GetMealsIDsByCartID(cart_id)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"filteredCards": FilteredMeals,
			"searchQuery":   childmealsquery,
			"cart_ID":       cart_id,
			"meals_cnt":     len(MealsIDs),
		})
	})

	r.POST("/home", func(c *gin.Context) {

		id := c.PostForm("add")
		log.Println(id)

		index, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant transform ind", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		a.repo.AddToCart(index)

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
				c.Error(err)
				return
			}
		}

		var cart_id int

		wrk_cart, err := a.repo.GetWorkingCart()
		log.Println(err)
		if len(wrk_cart) == 0 {
			new_cart, err := a.repo.CreateCart()
			cart_id = new_cart[0].ID
			log.Println(err)
		} else {
			cart_id = wrk_cart[0].ID
		}

		MealsIDs, err := a.repo.GetMealsIDsByCartID(cart_id)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"filteredCards": FilteredMeals,
			"searchQuery":   childmealsquery,
			"cart_ID":       cart_id,
			"meals_cnt":     len(MealsIDs),
		})

	})

	r.GET("/meal/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant get card by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		childmeal, err := a.repo.GetCardByID(index)
		if err != nil { // если не получилось
			log.Printf("cant get card by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "meal.html", gin.H{
			"title":     "Main website",
			"card_data": childmeal,
		})
	})

	r.GET("/cart/:id", func(c *gin.Context) {
		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get cart by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		MealsIDs, err := a.repo.GetMealsIDsByCartID(index)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")
			c.Error(err)
			return
		}

		MealsInCart := []ds.Meals{}
		for _, v := range MealsIDs {
			meal_tmp, err := a.repo.GetMealByID(v.ChildMealID)
			if err != nil {
				c.Error(err)
				return
			}
			MealsInCart = append(MealsInCart, meal_tmp[0])
			log.Println(v.ChildMealID)
		}

		c.HTML(http.StatusOK, "cart.html", gin.H{
			"title":     "Корзина",
			"CartMeals": MealsInCart,
			"CartID":    index,
		})
	})

	r.POST("/cart/:id", func(c *gin.Context) {

		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get cart by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		a.repo.DeleteCart(index)

		MealsIDs, err := a.repo.GetMealsIDsByCartID(index)
		if err != nil {
			log.Println("unable to get MealsIDsByCartID")
			c.Error(err)
			return
		}

		MealsInCart := []ds.Meals{}
		for _, v := range MealsIDs {
			meal_tmp, err := a.repo.GetMealByID(v.ChildMealID)
			if err != nil {
				c.Error(err)
				return
			}
			MealsInCart = append(MealsInCart, meal_tmp[0])
			log.Println(v.ChildMealID)
		}

		c.HTML(http.StatusOK, "cart.html", gin.H{
			"title":     "Корзина",
			"CartMeals": MealsInCart,
			"CartID":    index,
		})

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
