package crawler

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/internal/common"
	robotstxtutil "github.com/dashbikash/vidura-sense/internal/crawler/robotstxt-util"
)

var (
	log    = common.GetLogger()
	config = common.GetConfig()
)

type Crawler interface {
	Run(targetUrl string, checkRobot bool)
	OnSuccess(func(*http.Response))
	OnError(func(*http.Response))
	OnHtml(func(*goquery.Document))
	OnXml(func(*goquery.Document))
}

type BasicCrawler struct {
	Crawler
	client           *http.Client
	handlerOnSuccess func(*http.Response)
	handlerOnError   func(*http.Response)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func (cl *BasicCrawler) Run(targetUrl string, checkRobot bool) {
	if checkRobot {
		if allowed := robotstxtutil.IsAllowedUrl(targetUrl); !allowed {
			cl.handlerOnError(&http.Response{StatusCode: 403, Status: "403 Forbidden"})
			return
		}
	}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Error(err.Error())
	}
	req.Header.Set("User-Agent", config.Crawler.UserAgent)

	res, err := cl.client.Do(req)
	if err != nil {
		log.Error(err.Error())
	}

	if res.StatusCode == 200 {
		if cl.handlerOnSuccess != nil {
			cl.handlerOnSuccess(res)
		}

		if ok := strings.HasPrefix(res.Header.Get("Content-Type"), "text/html"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			doc.Url = res.Request.URL
			if cl.handlerOnHtml != nil {
				cl.handlerOnHtml(doc)
			}

		} else if ok := strings.HasPrefix(res.Header.Get("Content-Type"), "text/xml"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			if cl.handlerOnXml != nil {
				cl.handlerOnXml(doc)
			}

		}

	} else {
		if cl.handlerOnError != nil {
			cl.handlerOnError(res)
		}

	}

}
func (cl *BasicCrawler) OnSuccess(fn func(*http.Response)) {
	cl.handlerOnSuccess = fn
}
func (cl *BasicCrawler) OnError(fn func(*http.Response)) {
	cl.handlerOnError = fn
}
func (cl *BasicCrawler) OnHtml(fn func(*goquery.Document)) {
	cl.handlerOnHtml = fn
}
func (cl *BasicCrawler) OnXml(fn func(*goquery.Document)) {
	cl.handlerOnXml = fn
}

func New() Crawler {
	crawler := &BasicCrawler{client: &http.Client{}}
	return crawler
}

func randomProxy() *url.URL {
	i := rand.Intn(len(config.Crawler.Proxies) - 1)

	proxyURL, _ := url.Parse(config.Crawler.Proxies[i])
	log.Debug(proxyURL.String())
	return proxyURL
}
