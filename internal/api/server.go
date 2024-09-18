package api

import (
	"DevIntApp/internal/app/config"
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/dsn"
	"DevIntApp/internal/app/repository"
	"fmt"
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

	cards, err := a.repo.GetAllProducts()
	fmt.Println(err)

	r.SetFuncMap(template.FuncMap{
		"replaceNewline": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "/n", "<br>"))
		},
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Static("/css", "./resources")
	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {

		query := c.Query("query") // Получаем поисковый запрос из URL
		var find_prod []ds.Products
		if query == "" {
			find_prod = cards
		} else {
			find_prod, err = a.repo.GetProductByCardText(query)
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"filteredCards": find_prod,
			"searchQuery":   query,
		})
	})

	r.GET("/item_detailed/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil || index < 0 || index > len(cards) {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		card, err := a.repo.GetProductByID(index)
		if err != nil { // если не получилось
			log.Printf("cant get card by id %v", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "item_detailed.html", gin.H{
			"title":     "Main website",
			"card_data": card,
		})
	})

	r.GET("/cart", func(c *gin.Context) {

		c.HTML(http.StatusOK, "cart.html", gin.H{
			"title":      "Корзина",
			"cards_data": cards,
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
