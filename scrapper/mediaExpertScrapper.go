package scrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"strings"
)

type Result struct {
	name string
	url  string
}

type Feature struct {
	name  string
	value string
}

type MediaExpertScrapper struct {
	config ScrapperConfig
}

type Product struct {
	Name        string
	Description string
	Features    []Feature
	Url         string
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

		log.Printf("Waiting for body tag on %s product page...", product.name)
		// Send HTTP GET request to website URL
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.Navigate(product.url),
			chromedp.Sleep(mediaScrapper.config.delay),
			chromedp.WaitVisible("body", chromedp.ByQuery),
		)
		log.Println("Body tag was found.")
		if err != nil {
			log.Fatal(err)
		}

		var table string
		var description string
		log.Printf("Waiting for specification table and content on %s product page...", product.name)
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.WaitVisible(".content > .description", chromedp.ByQuery),
			chromedp.OuterHTML(".specification > table", &table, chromedp.ByQuery),
			chromedp.OuterHTML(".content > .description", &description, chromedp.ByQuery),
		)
		log.Printf("Specification table and content was found on %s product page.", product.name)
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

		scrapingResult := Product{
			Name:        product.name,
			Description: productDescription,
			Features:    features,
			Url:         product.url,
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
			feature.name = strings.TrimSpace(col.Text())
		})
		row.Find("td").Each(func(j int, col *goquery.Selection) {
			feature.value = strings.TrimSpace(col.Text())
		})
		features = append(features, feature)
	})

	return features
}

func (mediaScrapper *MediaExpertScrapper) getTenFirstResults(domain string, doc *goquery.Document) []Result {

	var results []Result

	doc.Find(".name.is-section").Each(func(i int, row *goquery.Selection) {
		result := Result{}

		result.name = strings.TrimSpace(row.Text())

		row.Find("a").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				href, exists := row.Attr("href")
				if exists {
					result.url = domain + href
				} else {
					result.url = "#"
				}

			}
		})
		results = append(results, result)
	})

	return results[:10]
}
