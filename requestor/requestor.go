package requestor

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/crawler"
	"github.com/dashbikash/vidura-sense/provider"
)

var (
	log = provider.GetLogger()
)

func Request() {
	crawler := crawler.New()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnHtml(func(d *goquery.Document) {
		log.Info(d.Find("title").First().Text())
		d.Find("h1").Each(func(i int, s *goquery.Selection) {
			log.Info(s.Text())
		})
		d.Find(".h1").Each(func(i int, s *goquery.Selection) {
			log.Info(s.Text())
		})
	})
	crawler.Run("https://quotes.toscrape.com")
}
