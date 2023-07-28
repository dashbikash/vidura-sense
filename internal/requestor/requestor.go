package requestor

import (
	"fmt"
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
	log    = system.Logger
	config = system.Config
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
		var links []string

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
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

	var stime, ftime time.Time

	// stime = time.Now()
	// crawler.RunManyAsync([]string{"https://quotes.toscrape.com", "https://www.metalsucks.net/"})
	// ftime = time.Now()
	// log.Info(fmt.Sprintf("Time elapsed %f", ftime.Sub(stime).Seconds()))

	stime = time.Now()
	crawler.RunManyAsync([]string{"https://quotes.toscrape.com", "https://www.metalsucks.net/"})
	ftime = time.Now()
	log.Info(fmt.Sprintf("Time elapsed %f", ftime.Sub(stime).Seconds()))
}

func RequestDemo2() {
	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(e error) {
		log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		var links []string

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
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

	var stime, ftime time.Time

	stime = time.Now()
	crawler.RunManyAsync([]string{"https://quotes.toscrape.com", "https://www.metalsucks.net/"})
	ftime = time.Now()
	log.Info(fmt.Sprintf("Time elapsed %f", ftime.Sub(stime).Seconds()))
}
