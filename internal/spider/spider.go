package spider

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	robotstxtutil "github.com/dashbikash/vidura-sense/internal/spider/robotstxt-util"
	"github.com/dashbikash/vidura-sense/internal/system"
)

var (
	log    = system.Logger
	config = system.Config
)

/*
 Spider interface for any spider
*/

type ISpider interface {
	RunManyAsync(targetUrls []string)
	RunManyAsyncAwait(targetUrls []string)
	RunMany(targetUrls []string)
	RunOneAsync(targetUrl string)
	RunOneAsyncAwait(targetUrl string)
	RunOne(targetUrl string)
	OnSuccess(func(*http.Response))
	OnError(func(error))
	OnHtml(func(*goquery.Document))
	OnXml(func(*goquery.Document))
}

/*
Basic spider struct implimenting ISpider
*/
type Spider struct {
	ISpider
	httpClient       *http.Client
	cfg              *SpiderConfig
	ctx              context.Context
	handlerOnSuccess func(*http.Response)
	handlerOnError   func(error)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func NewSpider() *Spider {

	c := context.WithValue(context.TODO(), "ProxyIndex", 0)
	spider := &Spider{httpClient: &http.Client{}, cfg: DefaultConfig(), ctx: c}

	return spider
}
func NewWithConfig(cfg *SpiderConfig) *Spider {

	spider := &Spider{httpClient: &http.Client{}, cfg: cfg}
	return spider
}

func (spider *Spider) makeRequest(targetUrl string) {

	if !spider.cfg.IgnoreRobotTxt {
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

func (spider *Spider) RunMany(targetUrls []string) {

	for _, t := range targetUrls {
		spider.makeRequest(t)
	}
}
func (spider *Spider) RunManyAsync(targetUrls []string) {
	var waitGroup sync.WaitGroup
	for _, t := range targetUrls {
		waitGroup.Add(1)
		t := t
		go func() {
			defer waitGroup.Done()
			spider.makeRequest(t)
		}()
	}
	// close the channel in the background
	waitGroup.Wait()
}
func (spider *Spider) RunManyAsyncAwait(targetUrls []string) {

	for _, t := range targetUrls {
		t := t
		go func() {
			spider.makeRequest(t)
		}()
	}
}

func (spider *Spider) RunOne(targetUrl string) {

	spider.makeRequest(targetUrl)
}
func (spider *Spider) RunOneAsync(targetUrl string) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		spider.makeRequest(targetUrl)
	}()

	waitGroup.Wait()
}
func (spider *Spider) RunOneAsyncAwait(targetUrl string) {

	go func() {
		spider.makeRequest(targetUrl)
	}()
}

func (spider *Spider) RandomProxy() *url.URL {
	i := rand.Intn(len(spider.cfg.Proxies) - 1)

	proxyURL, _ := url.Parse(spider.cfg.Proxies[i])
	log.Debug(proxyURL.String())
	return proxyURL
}

func (spider *Spider) RoundRobinProxy() *url.URL {
	spider.cfg.Proxies = config.Crawler.Proxies
	idx := spider.ctx.Value("ProxyIndex").(int)

	proxyURL, _ := url.Parse(spider.cfg.Proxies[idx])
	log.Debug(proxyURL.String())
	idx++
	if idx >= len(spider.cfg.Proxies) {
		idx = 0
	}
	spider.ctx = context.WithValue(spider.ctx, "ProxyIndex", idx)
	return proxyURL
}
