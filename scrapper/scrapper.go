package scrapper

func Scrap(domain string, url string, simplifiedProductNames bool) ([]Product, error) {
	mediaExpertScrapper := MediaExpertScrapper{
		config: getScrapperConfig(),
	}

	products, err := mediaExpertScrapper.Scrap(domain, url, simplifiedProductNames)

	if err != nil {
		return nil, err
	}

	return products, nil
}
