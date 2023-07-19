package restapi

import (
	"github.com/dashbikash/vidura-sense/provider"
	"github.com/dashbikash/vidura-sense/restapi/handlers"
	"github.com/gin-gonic/gin"
)

var (
	config = provider.GetConfig()
	router = gin.Default()
)

func setupServer() {

	router.GET("/", handlers.Index)

}
func Start() {
	setupServer()
	gin.SetMode(config.Server.Mode)
	router.Run(config.Server.Host + ":" + config.Server.Port)
}
