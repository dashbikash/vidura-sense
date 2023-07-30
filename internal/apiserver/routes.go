package apiserver

import (
	"github.com/dashbikash/vidura-sense/internal/apiserver/apihandler"
)

func setRoutes() {

	router.GET("/", apihandler.Index)

}
