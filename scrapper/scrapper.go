package scrapper

func Scrap() []Product {
	mediaExpertScrapper := MediaExpertScrapper{
		config: getScrapperConfig(),
	}

	products := mediaExpertScrapper.Scrap("https://www.mediaexpert.pl", "/agd-male/sprzatanie/odkurzacze-pionowe")

	return products
}
