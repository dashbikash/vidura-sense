package crawler

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/provider"
)

var (
	log    = provider.GetLogger()
	config = provider.GetConfig()
)

type Crawler interface {
	Run(targetUrl string)
	OnSuccess(func(*http.Response))
	OnError(func(*http.Response))
	OnHtml(func(*goquery.Document))
	OnXml(func(*goquery.Document))
}

type BasicCrawler struct {
	Crawler
	client           *http.Client
	handlerOnsuccess func(*http.Response)
	handlerOnError   func(*http.Response)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func (cl *BasicCrawler) Run(targetUrl string) {
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Error(err.Error())
	}
	req.Header.Set("User-Agent", config.Crawler.UserAgent)

	res, err := cl.client.Do(req)
	if err != nil {
		log.Error(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		cl.handlerOnsuccess(res)

		if _, ok := strings.CutPrefix(res.Header.Get("Content-Type"), "text/html"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			cl.handlerOnHtml(doc)
		}
		if _, ok := strings.CutPrefix(res.Header.Get("Content-Type"), "text/xml"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			cl.handlerOnXml(doc)
		}

	} else {
		cl.handlerOnError(res)
	}

}
func (cl *BasicCrawler) OnSuccess(fn func(*http.Response)) {
	cl.handlerOnsuccess = fn
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
