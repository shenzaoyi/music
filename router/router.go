package router

import (
	"Music/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine) {
	musicGroup := e.Group("/music/v1")
	{
		musicGroup.GET("/search", controller.SearchMusic)
		musicGroup.GET("/play", controller.PlayMusic)
		musicGroup.POST("/album", controller.GetAlbumMusics)
		musicGroup.GET("list", controller.GetAlbumList)
		//musicGroup.POST("/upload", controller.UploadMusic)
	}
}
