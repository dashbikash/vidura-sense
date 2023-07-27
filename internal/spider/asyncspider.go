package spider

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	robotstxtutil "github.com/dashbikash/vidura-sense/internal/spider/robotstxt-util"
)

type IAsyncSpider interface {
	Run(targetUrl []string, checkRobot bool)
	OnSuccess(func(*http.Response))
	OnError(func(error))
	OnHtml(func(*goquery.Document))
	OnXml(func(*goquery.Document))
}

type AsyncSpider struct {
	IAsyncSpider
	httpClient       *http.Client
	waitGroup        sync.WaitGroup
	handlerOnSuccess func(*http.Response)
	handlerOnError   func(error)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func (spider *AsyncSpider) makeRequest(targetUrl string, checkRobot bool) {
	defer spider.waitGroup.Done()
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
	req, err := http.NewRequest("GET", httpUrl.String(), nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.Header.Set("User-Agent", config.Crawler.UserAgent)

	resp, err := spider.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if spider.handlerOnSuccess != nil {
			spider.handlerOnSuccess(resp)
		}

		if ok := strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			doc.Url = resp.Request.URL
			if spider.handlerOnHtml != nil {
				spider.handlerOnHtml(doc)
			}

		} else if ok := strings.HasPrefix(resp.Header.Get("Content-Type"), "text/xml"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			if spider.handlerOnXml != nil {
				spider.handlerOnXml(doc)
			}

		}

	} else {
		if spider.handlerOnError != nil {
			spider.handlerOnError(errors.New(resp.Status))
		}

	}

}
func (spider *AsyncSpider) OnSuccess(fn func(*http.Response)) {
	spider.handlerOnSuccess = fn
}
func (spider *AsyncSpider) OnError(fn func(error)) {
	spider.handlerOnError = fn
}
func (spider *AsyncSpider) OnHtml(fn func(*goquery.Document)) {
	spider.handlerOnHtml = fn
}
func (spider *AsyncSpider) OnXml(fn func(*goquery.Document)) {
	spider.handlerOnXml = fn
}

func NewAsyncSpider() *AsyncSpider {

	spider := &AsyncSpider{httpClient: &http.Client{}}
	return spider
}

func (spider *AsyncSpider) Run(targetUrl []string, checkRobot bool) {

	for _, t := range targetUrl {
		spider.waitGroup.Add(1)
		go spider.makeRequest(t, checkRobot)
	}

	// close the channel in the background
	spider.waitGroup.Wait()

}

func randomProxy() *url.URL {
	i := rand.Intn(len(config.Crawler.Proxies) - 1)

	proxyURL, _ := url.Parse(config.Crawler.Proxies[i])
	log.Debug(proxyURL.String())
	return proxyURL
}
