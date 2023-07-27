package spider

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	robotstxtutil "github.com/dashbikash/vidura-sense/internal/spider/robotstxt-util"
	"github.com/dashbikash/vidura-sense/internal/system"
)

var (
	log    = system.GetLogger()
	config = system.GetConfig()
)

type ISpider interface {
	Run(targetUrl string, checkRobot bool)
	OnSuccess(func(*http.Response))
	OnError(func(error))
	OnHtml(func(*goquery.Document))
	OnXml(func(*goquery.Document))
}

type Spider struct {
	ISpider
	httpClient       *http.Client
	httpRequest      *http.Request
	handlerOnSuccess func(*http.Response)
	handlerOnError   func(error)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func (spider *Spider) Run(targetUrl string, checkRobot bool) {
	if checkRobot {
		if allowed := robotstxtutil.IsAllowedUrl(targetUrl); !allowed {
			spider.handlerOnError(errors.New("not allowed: robots.txt validation failed"))
			return
		}
	}
	httpUrl, err := url.Parse(targetUrl)
	if err != nil {
		spider.handlerOnError(errors.New("invalid url"))
		return
	}
	spider.httpRequest.URL = httpUrl
	spider.httpRequest.Host = httpUrl.Host
	res, err := spider.httpClient.Do(spider.httpRequest)

	if err != nil {
		log.Error(err.Error())
		return
	}

	if res.StatusCode == 200 {
		if spider.handlerOnSuccess != nil {
			spider.handlerOnSuccess(res)
		}

		if ok := strings.HasPrefix(res.Header.Get("Content-Type"), "text/html"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			doc.Url = res.Request.URL
			if spider.handlerOnHtml != nil {
				spider.handlerOnHtml(doc)
			}

		} else if ok := strings.HasPrefix(res.Header.Get("Content-Type"), "text/xml"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			if spider.handlerOnXml != nil {
				spider.handlerOnXml(doc)
			}

		}

	} else {
		if spider.handlerOnError != nil {
			spider.handlerOnError(errors.New(res.Status))
		}

	}

}
func (spider *Spider) OnSuccess(fn func(*http.Response)) {
	spider.handlerOnSuccess = fn
}
func (spider *Spider) OnError(fn func(error)) {
	spider.handlerOnError = fn
}
func (spider *Spider) OnHtml(fn func(*goquery.Document)) {
	spider.handlerOnHtml = fn
}
func (spider *Spider) OnXml(fn func(*goquery.Document)) {
	spider.handlerOnXml = fn
}

func NewSpider() *Spider {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	req.Header.Set("User-Agent", config.Crawler.UserAgent)

	spider := &Spider{httpClient: &http.Client{},
		httpRequest: req}
	return spider
}
