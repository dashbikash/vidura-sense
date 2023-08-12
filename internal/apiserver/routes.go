package apiserver

import (
	"github.com/dashbikash/vidura-sense/internal/apiserver/apihandler"
)

func setRoutes() {
	router.LoadHTMLGlob("internal/apiserver/templates/*")

	router.GET("/", apihandler.Index)

	urlRoute := router.Group("/url")
	{
		urlRoute.POST("/", apihandler.PostUrlAdd)
	}

	crawlRoutes := router.Group("/crawl")
	{
		crawlRoutes.POST("/new", apihandler.PostCrawlNew)
		crawlRoutes.POST("/url", apihandler.PostCrawlUrl)
		crawlRoutes.POST("/exclusive", apihandler.PostCrawlExclusive)
	}

}
