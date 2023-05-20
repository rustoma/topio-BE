package scrapper

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Result struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Feature struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MediaExpertScrapper struct {
	config ScrapperConfig
}

type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Features    []Feature `json:"features"`
	Url         string    `json:"url"`
	MainImage   string    `json:"main_image"`
	OrderWeight int       `json:"order_weight"`
}

func (mediaScrapper *MediaExpertScrapper) Scrap(domain string, url string, categories ...string) []Product {

	defer mediaScrapper.config.cancel()

	//Get first ten product names
	log.Println("Waiting for body tag on products list page...")
	err := chromedp.Run(mediaScrapper.config.ctx,
		chromedp.Navigate(domain+url),
		chromedp.Sleep(mediaScrapper.config.delay),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)
	log.Println("Body tag was found.")

	var body string

	err = chromedp.Run(mediaScrapper.config.ctx,
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("body", &body, chromedp.ByQuery),
	)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	products := mediaScrapper.getTenFirstResults(domain, doc)

	var scrapingResults []Product

	for _, product := range products {

		//var body string

		log.Printf("Waiting for body tag on %s product page...", product.Name)
		// Send HTTP GET request to website URL
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.Navigate(product.Url),
			chromedp.Sleep(mediaScrapper.config.delay),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			//chromedp.OuterHTML("body", &body, chromedp.ByQuery),
		)
		log.Println("Body tag was found.")

		//log.Println(body)

		if err != nil {
			log.Fatal(err)
		}

		var table string
		var description string
		var image string
		log.Printf("Waiting for specification table and content on %s product page...", product.Name)
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.WaitVisible(".item.description", chromedp.ByQuery),
			chromedp.OuterHTML(".specification > table", &table, chromedp.ByQuery),
			chromedp.OuterHTML(".item.description", &description, chromedp.ByQuery),
			chromedp.OuterHTML(".gallery-slide-image", &image, chromedp.ByQuery),
		)
		log.Printf("Specification table and content was found on %s product page.", product.Name)
		if err != nil {
			log.Fatal(err)
		}

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(table))
		if err != nil {
			log.Fatal(err)
		}

		features := mediaScrapper.getSpecifications(doc)

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(description))
		if err != nil {
			log.Fatal(err)
		}

		productDescription := mediaScrapper.getDescription(doc)

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(image))
		if err != nil {
			log.Fatal(err)
		}

		productImage := mediaScrapper.getProductImage(doc)

		scrapingResult := Product{
			Name:        product.Name,
			Description: productDescription,
			Features:    features,
			Url:         product.Url,
			MainImage:   productImage,
		}

		scrapingResults = append(scrapingResults, scrapingResult)
	}

	// Close the context and browser
	chromedp.Cancel(mediaScrapper.config.ctx)

	return scrapingResults
}

func (mediaScrapper *MediaExpertScrapper) getDescription(doc *goquery.Document) string {
	var result string

	doc.Each(func(i int, s *goquery.Selection) {

		s.Children().Each(func(i int, child *goquery.Selection) {
			text := strings.TrimSpace(child.Text())

			// Define the regular expression pattern to match consecutive whitespaces
			pattern := `\s+`

			// Compile the regular expression pattern
			reg := regexp.MustCompile(pattern)

			// Replace multiple whitespaces with a single space
			res := reg.ReplaceAllString(text, " ") + "\n\n"

			result += res + "\n\n"
		})
	})

	return result
}

func (mediaScrapper *MediaExpertScrapper) getSpecifications(doc *goquery.Document) []Feature {
	var features []Feature

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		feature := Feature{}
		row.Find("th").Each(func(j int, col *goquery.Selection) {
			feature.Name = strings.TrimSpace(col.Text())
		})
		row.Find("td").Each(func(j int, col *goquery.Selection) {
			feature.Value = strings.TrimSpace(col.Text())
		})
		features = append(features, feature)
	})

	return features
}

func (mediaScrapper *MediaExpertScrapper) getProductImage(doc *goquery.Document) string {

	var imageUrl string

	image := doc.Find(".image-magnifier > .spark-image").First()

	url, ok := image.Attr("src")
	if ok {
		imageUrl = url
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("Failed to fetch image:", err)

	}
	defer resp.Body.Close()

	// Read the image data into memory
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read image data:", err)
	}

	dataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData)

	return dataURI

}

func (mediaScrapper *MediaExpertScrapper) getTenFirstResults(domain string, doc *goquery.Document) []Result {

	var results []Result

	doc.Find(".name.is-section").Each(func(i int, row *goquery.Selection) {
		result := Result{}

		result.Name = strings.TrimSpace(row.Text())

		row.Find("a").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				href, exists := row.Attr("href")
				if exists {
					result.Url = domain + href
				} else {
					result.Url = "#"
				}

			}
		})
		results = append(results, result)
	})

	return results[:2]
}
