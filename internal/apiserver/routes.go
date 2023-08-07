package apiserver

import (
	"github.com/dashbikash/vidura-sense/internal/apiserver/apihandler"
)

func setRoutes() {

	router.GET("/", apihandler.Index)
	router.POST("/seedurl", apihandler.SeedUrl)
	router.POST("/crawl", apihandler.Crawl)

}
