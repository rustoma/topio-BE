package scrapper

func Scrap(domain string, url string) []Product {
	mediaExpertScrapper := MediaExpertScrapper{
		config: getScrapperConfig(),
	}

	products := mediaExpertScrapper.Scrap(domain, url)

	return products
}
