package restapi

import (
	"github.com/dashbikash/vidura-sense/internal/common"
	"github.com/dashbikash/vidura-sense/internal/restapi/handlers"
	"github.com/gin-gonic/gin"
)

var (
	config = common.GetConfig()
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
