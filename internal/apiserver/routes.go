package apiserver

import (
	"github.com/dashbikash/vidura-sense/internal/apiserver/apihandler"
)

func setRoutes() {
	router.LoadHTMLGlob("internal/apiserver/templates/*")

	router.GET("/", apihandler.Index)

	router.POST("/seedurl", apihandler.PostSeedUrl)

	crawl := router.Group("/crawl")
	{
		crawl.POST("/new", apihandler.PostCrawl)
		crawl.POST("/url", apihandler.PostCrawlUrl)
		crawl.POST("/exclusive", apihandler.PostCrawlUrl)
	}

}
