package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/magesh1/go-mongo/service"
)

func Routes(routes *gin.Engine) {
	url := routes.Group("/url")
	{
		url.POST("/short", service.ShortenUrl)
		url.GET("/redirect/:code", service.Redirect)
	}
}
