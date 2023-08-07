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
	AddUrlFilter(string, func(string) bool)
}

/*
Basic spider struct implimenting ISpider
*/
type Spider struct {
	ISpider
	httpClient       *http.Client
	cfg              *SpiderConfig
	ctx              context.Context
	urlFilters       map[string]func(string) bool
	handlerOnSuccess func(*http.Response)
	handlerOnError   func(error)
	handlerOnHtml    func(*goquery.Document)
	handlerOnXml     func(*goquery.Document)
}

func NewSpider() *Spider {

	c := context.WithValue(context.TODO(), "ProxyIndex", 0)
	spider := &Spider{httpClient: &http.Client{},
		cfg:        DefaultConfig(),
		ctx:        c,
		urlFilters: make(map[string]func(string) bool)}

	if !spider.cfg.IgnoreRobotTxt {
		spider.AddUrlFilter("RobotTxtValidation", func(targetUrl string) bool {
			return robotstxtutil.IsAllowedUrl(targetUrl)
		})
	}

	return spider
}
func NewWithConfig(cfg *SpiderConfig) *Spider {

	spider := &Spider{httpClient: &http.Client{},
		cfg:        cfg,
		urlFilters: make(map[string]func(string) bool)}
	return spider
}

func (spider *Spider) makeRequest(targetUrl string) {

	httpUrl, err := url.Parse(targetUrl)
	if err != nil {
		spider.handlerOnError(errors.New("invalid url"))
		return
	}
	httpReq, err := http.NewRequest("GET", httpUrl.String(), nil)
	if err != nil {
		system.Log.Error(err.Error())
		return
	}
	httpReq.Header.Set("User-Agent", system.Config.Crawler.UserAgent)

	resp, err := spider.httpClient.Do(httpReq)

	if err != nil {
		system.Log.Error(err.Error())
		return
	}
	defer func() {
		resp.Body.Close()
		resp = nil
		httpReq = nil
		httpUrl = nil
	}()

	if resp.StatusCode == 200 {
		if spider.handlerOnSuccess != nil {
			spider.handlerOnSuccess(resp)
		}

		if ok := strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			defer func() {
				doc = nil
			}()
			if err != nil {
				system.Log.Fatal(err.Error())
			}
			doc.Url = resp.Request.URL
			if spider.handlerOnHtml != nil {
				spider.handlerOnHtml(doc)
			}

		} else if ok := strings.HasPrefix(resp.Header.Get("Content-Type"), "text/xml"); ok {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			defer func() {
				doc = nil
			}()
			if err != nil {
				system.Log.Fatal(err.Error())
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

func (spider *Spider) AddUrlFilter(filterId string, fn func(targetUrl string) bool) {
	spider.urlFilters[filterId] = fn
}
func (spider *Spider) applyUrlFilters(targetUrl string) bool {
	for filterId, filterFn := range spider.urlFilters {

		if !filterFn(targetUrl) {
			system.Log.Info("Url filter failed at " + filterId)
			return false
		}
	}
	return true
}

func (spider *Spider) RunMany(targetUrls []string) {

	for _, t := range targetUrls {
		if spider.applyUrlFilters(t) {
			spider.makeRequest(t)
		}
	}
}
func (spider *Spider) RunManyAsyncAwait(targetUrls []string) {
	var waitGroup sync.WaitGroup
	for _, t := range targetUrls {
		waitGroup.Add(1)
		t := t
		go func() {
			defer waitGroup.Done()
			if spider.applyUrlFilters(t) {
				spider.makeRequest(t)
			}
		}()
	}
	// close the channel in the background
	waitGroup.Wait()
}
func (spider *Spider) RunManyAsync(targetUrls []string) {

	for _, t := range targetUrls {
		t := t
		go func() {
			if spider.applyUrlFilters(t) {
				spider.makeRequest(t)
			}
		}()
	}
}

func (spider *Spider) RunOne(targetUrl string) {

	if spider.applyUrlFilters(targetUrl) {
		spider.makeRequest(targetUrl)
	}
}
func (spider *Spider) RunOneAsyncAwait(targetUrl string) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		if spider.applyUrlFilters(targetUrl) {
			spider.makeRequest(targetUrl)
		}
	}()

	waitGroup.Wait()
}
func (spider *Spider) RunOneAsync(targetUrl string) {

	go func() {
		if spider.applyUrlFilters(targetUrl) {
			spider.makeRequest(targetUrl)
		}
	}()
}

func (spider *Spider) RandomProxy() *url.URL {
	i := rand.Intn(len(spider.cfg.Proxies) - 1)

	proxyURL, _ := url.Parse(spider.cfg.Proxies[i])
	system.Log.Debug(proxyURL.String())
	return proxyURL
}

func (spider *Spider) RoundRobinProxy() *url.URL {
	spider.cfg.Proxies = system.Config.Crawler.Proxies
	idx := spider.ctx.Value("ProxyIndex").(int)

	proxyURL, _ := url.Parse(spider.cfg.Proxies[idx])
	system.Log.Debug(proxyURL.String())
	idx++
	if idx >= len(spider.cfg.Proxies) {
		idx = 0
	}
	spider.ctx = context.WithValue(spider.ctx, "ProxyIndex", idx)
	return proxyURL
}
