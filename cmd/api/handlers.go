package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"topio/internal/models"
	"topio/scrapper"
)

type CreatePageRequest struct {
	Title        string `json:"title"`
	ProductName  string `json:"product_name"`
	Slug         string `json:"slug"`
	EntryDataUrl string `json:"entry_data_url"`
	ParentPage   *int   `json:"parent_page"`
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go topio up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) CreatePage(w http.ResponseWriter, r *http.Request) {
	var pageRequest CreatePageRequest

	err := app.readJSON(w, r, &pageRequest)

	questionsToGenerate := 2

	questions, err := app.AI.GenerateQuestionsForProduct(questionsToGenerate, pageRequest.ProductName)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("[CreatePage ERROR]: %v", err)
		}
		return
	}

	body, err := app.AI.GenerateBodyBasedOnQuestions(questions)

	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("[CreatePage ERROR]: %v", err)
		}
		return
	}

	page := models.Page{
		Title:        pageRequest.Title,
		ProductName:  pageRequest.ProductName,
		Body:         body,
		Slug:         pageRequest.Slug,
		EntryDataUrl: pageRequest.EntryDataUrl,
		ParentPage:   pageRequest.ParentPage,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	id, err := app.DB.InsertPage(page)

	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("[CreatePage ERROR]: %v", err)
		}
		return
	}

	_ = app.writeJSON(w, http.StatusOK, id)

}

type CreateProductsRequest struct {
	URL         string `json:"url"`
	RelatedPage int    `json:"related_page"`
}

func (app *application) CreateProducts(w http.ResponseWriter, r *http.Request) {

	var pageRequest CreateProductsRequest

	err := app.readJSON(w, r, &pageRequest)

	// Parse the URL
	parsedURL, err := url.Parse(pageRequest.URL)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		return
	}

	domain := parsedURL.Hostname()
	path := parsedURL.Path
	query := parsedURL.RawQuery

	results := scrapper.Scrap("https://"+domain, path+"?"+query)

	productsWithDescription := app.AI.GenerateProductsWithDescription(results)
	sortedProducts, err := app.AI.SortProductsFromBestToWorse(productsWithDescription)

	if err != nil {
		fmt.Println(err)
	}

	for index, product := range sortedProducts {

		product.RelatedPage = pageRequest.RelatedPage
		product.OrderWeight = index

		id, err := app.DB.InsertProduct(product)

		if err != nil {
			println(err)
		}

		fmt.Printf("Product %s created with id %d", product.Name, id)

	}

	_ = app.writeJSON(w, http.StatusOK, results)
}
