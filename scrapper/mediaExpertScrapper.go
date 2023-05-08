package scrapper

import (
	"fmt"
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

func (mediaScrapper *MediaExpertScrapper) Scrap(domain string, url string, categories ...string) {

	defer mediaScrapper.config.cancel()

	//Get first ten product names
	err := chromedp.Run(mediaScrapper.config.ctx,
		chromedp.Navigate(domain+url),
		chromedp.Sleep(mediaScrapper.config.delay),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)

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

	for _, product := range products {

		// Send HTTP GET request to website URL
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.Navigate(product.url),
			chromedp.Sleep(mediaScrapper.config.delay),
			chromedp.WaitVisible("body", chromedp.ByQuery),
		)

		if err != nil {
			log.Fatal(err)
		}

		var table string
		var description string
		err = chromedp.Run(mediaScrapper.config.ctx,
			chromedp.WaitVisible(".specification > table", chromedp.ByQuery),
			chromedp.WaitVisible(".content.content-shadow", chromedp.ByQuery),
			chromedp.OuterHTML(".specification > table", &table, chromedp.ByQuery),
			chromedp.OuterHTML(".content.content-shadow", &description, chromedp.ByQuery),
		)

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

		textForAI := mediaScrapper.getDescription(doc)
		fmt.Printf("%+v\n", features)
		fmt.Println(textForAI)

	}

	// Close the context and browser
	chromedp.Cancel(mediaScrapper.config.ctx)

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
			feature.name = col.Text()
		})
		row.Find("td").Each(func(j int, col *goquery.Selection) {
			feature.value = col.Text()
		})
		features = append(features, feature)
	})

	return features
}

func (mediaScrapper *MediaExpertScrapper) getTenFirstResults(domain string, doc *goquery.Document) []Result {

	var results []Result

	doc.Find(".name.is-section").Each(func(i int, row *goquery.Selection) {
		result := Result{}

		result.name = row.Text()

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
