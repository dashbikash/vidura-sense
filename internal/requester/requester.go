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

func SimpleRequest() {
	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		system.Log.Info(r.Status + " " + r.Request.URL.String())
	})
	crawler.OnError(func(e error) {
		system.Log.Error(e.Error())
	})
	crawler.OnHtml(func(d *goquery.Document) {
		links := []string{}

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
			}
			// go func() {
			// 	if strings.HasPrefix(href, "http") {
			// 		crawler.RunOne(href)
			// 	}
			// }()
		})
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
		//system.Log.Debug(html.UpdatedOn.Format("yyyy-MM-dd"))

	})
	crawler.OnXml(func(d *goquery.Document) {
		system.Log.Info(d.Find("rss").First().Text())
	})

	crawler.RunManyAsyncAwait([]string{"https://quotes.toscrape.com"})
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

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
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
			html.StoreHtml()
		}()

	})
	crawler.OnXml(func(d *goquery.Document) {
		system.Log.Info(d.Find("rss").First().Text())
	})

	crawler.RunManyAsync([]string{"https://quotes.toscrape.com"})
}
