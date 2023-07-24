package requestor

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/dashbikash/vidura-sense/internal/common"
	"github.com/dashbikash/vidura-sense/internal/crawler"
	redisstore "github.com/dashbikash/vidura-sense/internal/datastore/redis-store"
)

var (
	log = common.GetLogger()
)

func RequestDemo() {
	crawler := crawler.New()
	crawler.OnSuccess(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnError(func(r *http.Response) {
		log.Info(r.Status)
	})
	crawler.OnHtml(func(d *goquery.Document) {

		log.Info(d.Find("title").First().Text())

		re := regexp.MustCompile(`\$\{(.*?)\}`)
		match := re.FindStringSubmatch(d.Find("body").First().Text())
		log.Println(match[0])

	})
	crawler.OnXml(func(d *goquery.Document) {
		log.Info(d.Find("rss").First().Text())
	})

	crawler.Run("https://quotes.toscrape.com")
	//crawler.Run("http://rss.cnn.com/rss/edition.rss")
}

func GetRobots() {
	crawler := crawler.New()
	crawler.OnSuccess(func(r *http.Response) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		redisstore.SetRobotsTxt(r.Request.Host, string(body))
	})
	crawler.OnError(func(r *http.Response) {
		log.Info(r.Status)
	})

	crawler.Run("https://stackoverflow.com/robots.txt")
	//crawler.Run("http://rss.cnn.com/rss/edition.rss")
}
