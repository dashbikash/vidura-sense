package apiserver

import (
	"github.com/dashbikash/vidura-sense/internal/apiserver/apihandler"
)

func setRoutes() {

	router.GET("/", apihandler.Index)
	router.POST("/seedurl", apihandler.PostSeedUrl)
	router.POST("/crawl", apihandler.PostCrawl)
	router.POST("/crawlurl", apihandler.PostCrawlUrl)

}
