package app

import (
	"log"
	"net/http"
	"parser/internal/config"
	"parser/internal/model"
	"parser/internal/repository"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	bookRepository *repository.Repository
}

func ParseBooks(url string) ([]model.Book, error) {
	res, err := http.Get(url)
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
	return books, nil
}
