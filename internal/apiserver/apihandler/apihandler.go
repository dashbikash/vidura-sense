package apihandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dashbikash/vidura-sense/internal/data/entity"
	"github.com/dashbikash/vidura-sense/internal/requester"
	urlfrontier "github.com/dashbikash/vidura-sense/internal/url-frontier"
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}
func PostUrlAdd(ctx *gin.Context) {
	var blankPages []interface{}
	var urls []string
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if ctx.ShouldBind(&urls) == nil {
		for _, targetUrl := range urls {
			blankPages = append(blankPages, entity.NewBlankHtmlPage(targetUrl))
		}
		entity.HtmlPageCreateBlankEntries(&blankPages)
	}

	ctx.String(200, fmt.Sprintf("%d urls seeded.", len(blankPages)))
}
func PostCrawlNew(ctx *gin.Context) {

	limit, err := strconv.ParseInt(ctx.DefaultQuery("lim", "30"), 0, 0)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	urls := urlfrontier.GetNewUrls(int(limit))
	requester.BroadRequest(urls)
	ctx.String(200, "Done")
}

func PostCrawlUrl(ctx *gin.Context) {

	var urls []string

	if err := ctx.ShouldBind(&urls); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	requester.BroadRequest(urls)
	ctx.String(200, "Done")
}
func PostCrawlExclusive(ctx *gin.Context) {

	limit, err := strconv.ParseInt(ctx.DefaultQuery("lim", "30"), 0, 0)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var domains []string

	if err := ctx.ShouldBind(&domains); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	urls := urlfrontier.GetExclusiveDomainNewUrls(int(limit), domains)
	requester.ExclusiveDomainRequest(urls, domains)
	ctx.String(200, "Done")
}
