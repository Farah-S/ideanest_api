package routes

import (
	"github.com/example/golang-test/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func SignUpRoute(router *gin.Engine)  {
    router.POST("/signup", controllers.CreateUser()) //add this
}

func IndexRoute(router *gin.Engine, port string){
    router.LoadHTMLGlob("pkg/pages/*")
	router.GET("/api", controllers.DisplayIndexPage())
	router.Run(":" + port)
}