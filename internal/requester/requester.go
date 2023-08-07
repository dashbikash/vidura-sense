package requester

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/internal/data/entity"
	"github.com/dashbikash/vidura-sense/internal/data/processor"
	"github.com/dashbikash/vidura-sense/internal/spider"
	"github.com/dashbikash/vidura-sense/internal/system"
)

func SimpleRequest(targetUrl []string) {

	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		system.Log.Info(r.Status + " " + r.Request.URL.String())
	})
	crawler.OnError(func(e error) {
		system.Log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		blankPages := make([]interface{}, 0)

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				if strings.HasPrefix(href, "http") {
					blankPages = append(blankPages, entity.BlankHtmlPage(strings.Trim(href, "/")))
				}
			}

		})
		html := &entity.HtmlPage{
			URL:       strings.Trim(d.Url.String(), "/"),
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: time.Now().Local(),
		}

		defer func() {
			blankPages, d, html = nil, nil, nil
		}()

		html.StoreUpdated()
		if len(blankPages) > 0 {
			entity.HtmlPageCreateBlankEntries(&blankPages)
		}

	})
	crawler.OnXml(func(d *goquery.Document) {
		system.Log.Info(d.Find("rss").First().Text())
	})
	// crawler.AddUrlFilter("TimeIntervalValidation", func(targetUrl string) bool {
	// 	urlMeta := &entity.UrlMeta{}
	// 	entity.GetUrlMeta("quotes.toscrape.com", urlMeta)
	// 	return time.Since(urlMeta.LastUpdated).Minutes() > 60
	// })

	crawler.RunMany(targetUrl)
}

func RecurssiveRequest() {
	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		system.Log.Info(r.Status + " " + r.Request.URL.String())
	})
	crawler.OnError(func(e error) {
		system.Log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		links := []string{}
		blankPages := make([]interface{}, 0)

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				blankPages = append(blankPages, entity.BlankHtmlPage(href))
			}

			go func() {
				href = regexp.MustCompile(`(.+?)(\#[^.]*$|$)`).ReplaceAllString(href, "${1}")
				purl, e := url.Parse(href)
				if e != nil {
					system.Log.Error(e.Error())
				}
				if purl.Scheme == "" || purl.Hostname() == d.Url.Hostname() {

					system.Log.Info("Visiting " + "http://" + d.Url.Hostname() + href)
					crawler.RunOne("http://" + d.Url.Hostname() + href)
				}
			}()

		})
		go func() {
			html := &entity.HtmlPage{
				URL:       strings.Trim(d.Url.String(), "/"),
				Scheme:    d.Url.Scheme,
				Host:      d.Url.Host,
				Title:     d.Find("title").First().Text(),
				Links:     links,
				Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
				UpdatedOn: time.Now().Local(),
			}
			html.StoreUpdated()
		}()

	})
	crawler.OnXml(func(d *goquery.Document) {
		system.Log.Info(d.Find("rss").First().Text())
	})

	crawler.RunManyAsync([]string{"https://quotes.toscrape.com"})
}
