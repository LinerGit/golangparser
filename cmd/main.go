package main

// @title Parser API
// @version 1.0
// @description Chi swagger
// @host localhost:3000
// @BasePath /

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"parser/internal/app"
	"parser/internal/config"
	"parser/internal/database"
	"parser/internal/logger"
	"parser/internal/model"
	"parser/internal/repository"
	"strconv"

	_ "parser/docs"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
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
		log.Fatalf("failed to load config %v", err)
	}

	logger := logger.New(cfg)
	db, err := database.NewDB()
	repo := repository.NewRepository(db, &logger)

	_ = repo

	r := chi.NewRouter()

	r.Get("/ping", PingExample)
	//r.Get("/getbook/{id}", http)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/job/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(idParam)

		book, err := repository.GetBookById(context.Background(), uint(id))
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		json.NewEncoder(w).Encode(book)
	})

	http.ListenAndServe(":3000", r)

	db.AutoMigrate(model.Book{})

	books, err := app.ParseBooks("https://books.toscrape.com/")

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
