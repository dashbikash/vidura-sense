package requester

import (
	"net/http"
	"net/url"
	"slices"
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

func BroadRequest(targetUrl []string) {

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

			if href, ok := s.Attr("href"); ok {
				if hrefParse, err := url.Parse(href); err == nil {
					if hrefParse.Host == "" {
						hrefParse.Host = d.Url.Host
						hrefParse.OmitHost = false
					}
					if hrefParse.Scheme == "" {
						hrefParse.Scheme = d.Url.Scheme
					}
					if strings.HasPrefix(hrefParse.Scheme, "http") {
						blankPages = append(blankPages, entity.NewBlankHtmlPage(strings.Trim(hrefParse.String(), "/")))
					}
				}
			}
		})
		d.Find("meta").First()

		html := &entity.HtmlPage{
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: primitive.DateTime(time.Now().UnixMilli()),
		}

		d.Url.Scheme = ""

		if metaContent, ok := d.Find("meta[name='description' i]").First().Attr("content"); ok {
			html.Meta.Description = metaContent
		}
		if metaContent, ok := d.Find("meta[name='author' i]").First().Attr("content"); ok {
			html.Meta.Author = metaContent
		}

		if metaContent, ok := d.Find("meta[name='keywords' i]").First().Attr("content"); ok {
			html.Meta.Keywords = metaContent
		}
		if metaContent, ok := d.Find("meta[name='viewport' i]").First().Attr("content"); ok {
			html.Meta.Viewport = metaContent
		}
		if metaContent, ok := d.Find("meta[charset]").First().Attr("charset"); ok {
			html.Meta.Charset = metaContent
		}

		html.URL = strings.Trim(d.Url.String(), "/")

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

func ExclusiveDomainRequest(targetUrl []string, filterDomains []string) {

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

			if href, ok := s.Attr("href"); ok {
				if hrefParse, err := url.Parse(href); err == nil {
					if hrefParse.Host == "" {
						hrefParse.Host = d.Url.Host
						hrefParse.OmitHost = false
					}
					if hrefParse.Scheme == "" {
						hrefParse.Scheme = d.Url.Scheme
					}
					if strings.HasPrefix(hrefParse.Scheme, "http") {
						blankPages = append(blankPages, entity.NewBlankHtmlPage(strings.Trim(hrefParse.String(), "/")))
					}
				}
			}
		})

		html := &entity.HtmlPage{
			Scheme:    d.Url.Scheme,
			Host:      d.Url.Host,
			Title:     d.Find("title").First().Text(),
			Body:      processor.NewTextProcessor(d.Find("body").Text()).SanitizeHtml().SanitizeText().String(),
			UpdatedOn: primitive.DateTime(time.Now().UnixMilli()),
		}

		d.Url.Scheme = ""

		if metaContent, ok := d.Find("meta[name='description' i]").First().Attr("content"); ok {
			html.Meta.Description = metaContent
		}
		if metaContent, ok := d.Find("meta[name='author' i]").First().Attr("content"); ok {
			html.Meta.Author = metaContent
		}

		if metaContent, ok := d.Find("meta[name='keywords' i]").First().Attr("content"); ok {
			html.Meta.Keywords = metaContent
		}
		if metaContent, ok := d.Find("meta[name='viewport' i]").First().Attr("content"); ok {
			html.Meta.Viewport = metaContent
		}
		if metaContent, ok := d.Find("meta[charset]").First().Attr("charset"); ok {
			html.Meta.Charset = metaContent
		}

		html.URL = strings.Trim(d.Url.String(), "/")

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

	crawler.AddUrlFilter("ExclusiveDomainFilter", func(targetUrl string) bool {
		if purl, err := url.Parse(targetUrl); err == nil {
			return slices.Contains(filterDomains, purl.Host)
		}

		return false
	})

	crawler.RunManyAsync(targetUrl)
}
