package scrapper

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/mileusna/useragent"
	"math/rand"
	"time"
)

type ScrapperConfig struct {
	ctx    context.Context
	delay  time.Duration
	cancel context.CancelFunc
}

func getScrapperConfig() ScrapperConfig {

	ctx, cancel := setupChromeDpContext()

	config := ScrapperConfig{
		ctx:    ctx,
		delay:  getRandomDelay(),
		cancel: cancel,
	}
	return config
}

func setupChromeDpContext() (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), getChromeDpOptions()...)

	ctx, cancel = chromedp.NewContext(ctx)
	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)

	return ctx, cancel
}

func getRandomDelay() time.Duration {
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(rand.Intn(5)) * time.Second
	return delay
}

func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
	}
	randomIndex := rand.Intn(len(userAgents))

	return userAgents[randomIndex]
}

func getChromeDpOptions() []chromedp.ExecAllocatorOption {
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(useragent.Parse(getRandomUserAgent()).String),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-logging", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("mute-audio", true),
	}
	return options
}
