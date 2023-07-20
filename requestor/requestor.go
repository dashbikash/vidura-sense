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

func RequestDemo() {
	crawler := crawler.New()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnHtml(func(d *goquery.Document) {

		log.Info(d.Find("title").First().Text())
		d.Find("h1,h2,h3,h4,h5,h6,span,p,pre").Each(func(i int, s *goquery.Selection) {
			log.Debugf("%s: %s", goquery.NodeName(s), s.Text())
		})
	})
	crawler.OnXml(func(d *goquery.Document) {
		log.Info(d.Find("rss").First().Text())
	})

	crawler.Run("https://quotes.toscrape.com")
	//crawler.Run("http://rss.cnn.com/rss/edition.rss")
}
