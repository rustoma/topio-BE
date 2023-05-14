package main

import (
	"fmt"
	"net/http"
	"topio/scrapper"
)

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

	results := scrapper.Scrap()

	productsWithDescription := app.AI.GenerateProductsWithDescription(results)
	sortedProducts, err := app.AI.SortProductsFromBestToWorse(productsWithDescription)

	if err != nil {
		fmt.Println(err)
	}

	for _, product := range sortedProducts {
		println("Product name: ")
		println(product.Name)
		println("Product description: ")
		println(product.GeneratedDescription)
		println("")
		println("Product features: ")
		fmt.Printf("%v", product.Features)
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
