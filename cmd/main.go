package main

// @title Parser API
// @version 1.0
// @description Chi swagger
// @host localhost:3000
// @BasePath /

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"parser/internal/config"
	"parser/internal/logger"
	"parser/internal/model"
	"parser/internal/repository"
	"strings"

	_ "parser/docs"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ping godoc
// @Summary Проверка сервера
// @Description Возвращает pong
// @Tags health
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingExample(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config ", err)
	}

	logger := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
	if err != nil {
		logger.Error().Msgf("failed to conn to db")
		return
	}
	logger.Info().Msg("db connected")

	repo := repository.NewRepository(db, &logger)

	_ = repo

	r := chi.NewRouter()

	r.Get("/ping", PingExample)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	http.ListenAndServe(":3000", r)

	db.AutoMigrate(model.Book{})

	res, err := http.Get("https://books.toscrape.com/")
	if err != nil {
		log.Fatal("Failed to connect to the target page", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("HTTP Error %d: %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	var books []model.Book
	doc.Find("article.product_pod").Each(func(i int, p *goquery.Selection) {
		book := model.Book{}
		book.Name = p.Find("a").Text()
		book.Price = p.Find("p.price_color").Text()
		book.Stock = p.Find("p.instock.availability").Text()
		book.Stock = strings.TrimSpace(book.Stock)
		books = append(books, book)
	})

	db.Create(&books)

	dataj, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("books.json", dataj, 0644)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("books.csv")

	if err != nil {
		log.Fatal("Failed to create the output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"name",
		"price",
		"stock",
	}

	writer.Write(headers)

	for _, book := range books {

		record := []string{
			book.Name,
			book.Price,
			book.Stock,
		}

		writer.Write(record)
	}

}
