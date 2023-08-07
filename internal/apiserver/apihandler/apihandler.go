package apihandler

import (
	"fmt"

	"github.com/dashbikash/vidura-sense/internal/data/entity"
	"github.com/dashbikash/vidura-sense/internal/requester"
	urlfrontier "github.com/dashbikash/vidura-sense/internal/url-frontier"
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.String(200, "Welcome to Vidura Sense")
}
func SeedUrl(ctx *gin.Context) {
	var blankPages []interface{}
	var urls []string
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if ctx.ShouldBind(&urls) == nil {
		for _, targetUrl := range urls {
			blankPages = append(blankPages, entity.BlankHtmlPage(targetUrl))
		}
		entity.HtmlPageCreateBlankEntries(&blankPages)
	}

	ctx.String(200, fmt.Sprintf("%d urls seeded.", len(blankPages)))
}
func Crawl(ctx *gin.Context) {
	urls := urlfrontier.GetNewUrls(30)
	requester.SimpleRequest(urls)
	ctx.String(200, "Done")
}
