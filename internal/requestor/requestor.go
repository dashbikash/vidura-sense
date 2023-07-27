package requestor

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/internal/data/entity"
	"github.com/dashbikash/vidura-sense/internal/data/processor"
	"github.com/dashbikash/vidura-sense/internal/spider"
	"github.com/dashbikash/vidura-sense/internal/system"
)

var (
	log = system.GetLogger()
)

func RequestDemo1() {
	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(e error) {
		log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		var links []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}
		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, struct {
					Title string `json:"title"`
					URL   string `json:"url"`
				}{
					Title: s.Text(),
					URL:   href,
				})
			}

		})
		log.Info("Processing Html for " + d.Url.String())
		html := &entity.HtmlPage{
			URL:       strings.Trim(d.Url.String(), "/"),
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Links:     links,
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: time.Now().Local(),
		}
		html.StoreHtml()
	})
	crawler.OnXml(func(d *goquery.Document) {
		log.Info(d.Find("rss").First().Text())
	})

	crawler.Run("https://quotes.toscrape.com", true)
	// crawler.Run("https://quotes.toscrape.com/tag/love/", true)
	// crawler.Run("https://quotes.toscrape.com/tag/life/", true)
	// crawler.Run("https://quotes.toscrape.com/tag/books/", true)
	// crawler.Run("http://rss.cnn.com/rss/edition.rss")
}

func RequestDemo2() {
	crawler := spider.NewAsyncSpider()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(e error) {
		log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		var links []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}
		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, struct {
					Title string `json:"title"`
					URL   string `json:"url"`
				}{
					Title: s.Text(),
					URL:   href,
				})
			}

		})
		log.Info("Processing Html for " + d.Url.String())
		html := &entity.HtmlPage{
			URL:       strings.Trim(d.Url.String(), "/"),
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Links:     links,
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: time.Now().Local(),
		}
		html.StoreHtml()
	})
	crawler.OnXml(func(d *goquery.Document) {
		log.Info(d.Find("rss").First().Text())
	})

	crawler.Run([]string{"https://quotes.toscrape.com", "https://quotes.toscrape.com/tag/love/", "https://quotes.toscrape.com/tag/book/", "https://godoc.org", "https://www.packtpub.com", "https://kubernetes.io/"}, true)
	// crawler.Run("https://quotes.toscrape.com/tag/books/", true)
	// crawler.Run("http://rss.cnn.com/rss/edition.rss")
}
