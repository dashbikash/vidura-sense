package apihandler

import (
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.String(200, "Welcome to Vidura Sense")
}
