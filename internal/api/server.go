package api

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Card struct {
	CardText string
	ImageUrl string
	HrefID   int
}

type CardsDetailsStruct struct {
	CardDetail string
}

func CardsInfoFunc() []Card {
	CardsInfo := []Card{
		{"Агуша Молоко стерилизованное детское 3,2%, 900 мл",
			"http://localhost:9000/development-internet-applications/photo_2024-09-04_09-14-57.jpg",
			1},
		{"Смесь детская молочная HiPP Combiotic 2, с 6 месяцев, 600г",
			"http://localhost:9000/development-internet-applications/smes_detskaya.jpg",
			2},
		{"Сок ФрутоНяня Яблоко осветленный, c 4 месяцев , 0.2 л",
			"http://localhost:9000/development-internet-applications/sok_yabloko_detskii_1.jpg",
			3},
		{"Пюре ФрутоНяня из кабачков с молоком, с 6 месяцев, 80 г",
			"http://localhost:9000/development-internet-applications/kabachok-moloko.jpg",
			4}}
	return CardsInfo
}

func CardsDetailed() []CardsDetailsStruct {
	CardsDetails := []CardsDetailsStruct{
		{"Бренд: Агуша\nТип: молоко\nМинимальный возраст: 3 года\nСырье: коровье молоко\nОбработка молока: стерилизованное\nЖирность: 3.2 %\nУпаковка: тетра-пак\nСостав\nМолоко нормализованное\nБелки в 100 г: 3 г\nЖиры в 100 г: 3.2 г\nУглеводы в 100 г: 4.7 г\n"},
		{"Бренд: HiPP\nЛинейка: Combiotic\nСтупень: 2\nРекомендуемый возраст: с 6 месяцев\nФорма выпуска: Сухая адаптированная\nНе содержит: ГМО\nСодержит: пребиотики, пробиотики\nВитамины: С, ниацин, Е, пантотеновая кислота, В1, В6, К, А, фолиевая кислота, биотин, В2, D, В12, биотин, холин, инозит, таурин, L-карнитин\nМинералы: калий, кальций, фосфор, натрий, железо, магний, хлориды, цинк, медь, йод, марганец, селен\nУпаковка: картонная коробка\nСостав: витамин е,натрия цитрат,пальмовое масло,мальтодекстрин,рапсовое масло,витамин а,витамин в3,сухая молочная сыворотка,витамин с,пищевые волокна,лактулоза,витамин,витамин d,пантотеновая кислота,фолиевая кислота,марганца сульфат,молоко,витамин к,пробиотическая молочнокислая культура,железа лактат\n"},
		{"Бренд: ФрутоНяня\nТип: сок\nМинимальный возраст: 4 месяца\nВкус: яблоко\nОсобенности: гипоаллергенно\nНе содержит: ГМО\nТип упаковки: Tetra Pak\nОбъем упаковки: 0.2 л\nВес: 0.1 кг\nВес: 0.1 кг\nСостав: Сок из яблок\nИзготовлен из концентрированного сока.\n"},
		{"Бренд: ФрутоНяня\nТип: однокомпонентное\nКомпоненты: овощи\nМинимальный возраст: 6 месяцев\nОвощи: брокколи, кабачок\nФрукты и ягоды: абрикос\nМясо: телятина\nПтица: индейка\nСубпродукты: сердце\nЗлаки: смесь злаков\nНе содержит: консерванты\nДобавки: творог\nОсобенности: гипоаллергенно\nУпаковка: стеклянная банка\nВес: 80 г\n"},
	}
	return CardsDetails
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"replaceNewline": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "\n", "<br>"))
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
		cards := CardsInfoFunc()

		// Фильтрация карточек по запросу
		var filteredCards []Card
		for _, card := range cards {
			if strings.Contains(strings.ToLower(card.CardText), strings.ToLower(query)) {
				filteredCards = append(filteredCards, card)
			}
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":         "Заказы на молочную кухню",
			"cards_data":    CardsInfoFunc(),
			"filteredCards": filteredCards,
			"searchQuery":   query,
		})
	})

	r.GET("/item_detailed/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil || index < 0 || index > len(CardsInfoFunc()) {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "item_detailed.html", gin.H{
			"title":              "Main website",
			"card_data":          CardsInfoFunc()[index-1],
			"card_detailed_data": CardsDetailed()[index-1],
		})
	})

	r.GET("/cart", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cart.html", gin.H{
			"title":      "Корзина",
			"cards_data": CardsInfoFunc(),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
