package requestor

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/internal/common"
	"github.com/dashbikash/vidura-sense/internal/crawler"
	"github.com/dashbikash/vidura-sense/internal/dataprocessor"
	"github.com/dashbikash/vidura-sense/internal/datatype"
)

var (
	log = common.GetLogger()
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
		log.Info("Processing Html for " + d.Url.String())
		html := &datatype.HtmlPage{
			URL:       d.Url.String(),
			Scheme:    d.Url.Scheme,
			Title:     d.Find("title").First().Text(),
			Body:      dataprocessor.SanitizeText(dataprocessor.SanitizeHtml(d.Find("body").Text())),
			UpdatedOn: time.Now(),
		}
		dataprocessor.StoreHtml(html)
	})
	crawler.OnXml(func(d *goquery.Document) {
		log.Info(d.Find("rss").First().Text())
	})

	go crawler.Run("https://quotes.toscrape.com/", true)
	go crawler.Run("https://quotes.toscrape.com/tag/love/", true)
	go crawler.Run("https://quotes.toscrape.com/tag/life/", true)
	go crawler.Run("https://quotes.toscrape.com/tag/books/", true)
	//crawler.Run("http://rss.cnn.com/rss/edition.rss")
}
