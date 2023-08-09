package requester

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/avast/retry-go/v4"
	"github.com/dashbikash/vidura-sense/internal/data/entity"
	"github.com/dashbikash/vidura-sense/internal/data/processor"
	"github.com/dashbikash/vidura-sense/internal/spider"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SimpleRequest(targetUrl []string) {

	crawler := spider.NewSpider()
	crawler.OnSuccess(func(r *http.Response) {
		system.Log.Info(r.Status + " " + r.Request.URL.String())
	})
	crawler.OnError(func(s string, err error) {
		system.Log.Error(err.Error())

		if _, ok := err.(retry.Error); ok {

			entity.HtmlPageLockUpdate(s)
		}
	})
	crawler.OnHtml(func(d *goquery.Document) {
		blankPages := make([]interface{}, 0)

		d.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				if strings.HasPrefix(href, "http") {
					blankPages = append(blankPages, entity.NewBlankHtmlPage(strings.Trim(href, "/")))
				}
			}

		})
		html := &entity.HtmlPage{
			URL:       strings.Trim(d.Url.String(), "/"),
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: primitive.DateTime(time.Now().UnixMilli()),
		}

		defer func() {
			blankPages, d, html = nil, nil, nil
		}()

		html.Store()
		if len(blankPages) > 0 {
			entity.HtmlPageCreateBlankEntries(&blankPages)
		}

	})
	crawler.OnXml(func(d *goquery.Document) {
		system.Log.Info(d.Find("rss").First().Text())
	})

	crawler.RunManyAsync(targetUrl)
}
