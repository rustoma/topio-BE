package scrapper

func Scrap() {
	mediaExpertScrapper := MediaExpertScrapper{
		config: getScrapperConfig(),
	}

	mediaExpertScrapper.Scrap("https://www.mediaexpert.pl", "/agd-male/sprzatanie/odkurzacze-pionowe")
}
