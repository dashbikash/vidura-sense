package job

import (
	"github.com/dashbikash/vidura-sense/internal/requester"
	urlfrontier "github.com/dashbikash/vidura-sense/internal/url-frontier"
)

func CrawlNewHtmlPages() {

	requester.SimpleRequest(urlfrontier.GetNewUrls(30))
}
