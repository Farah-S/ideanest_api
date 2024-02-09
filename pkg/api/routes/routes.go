package routes

import (
	"net/http"

	// "github.com/example/golang-test/pkg/api/handlers"
	"github.com/example/golang-test/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func StartApp(router *gin.Engine)  {
	// r := gin.Default()
	// r.GET("/api", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	SignUpFormRoute(router)

	// r.GET("/success", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "success.html", nil)
	// })

	router.Run(":8080")
}

func SignUpFormRoute(router *gin.Engine)  {
    router.GET("/api/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	}) //add this
	SignUpRoute(router)
	router.GET("/success", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
}

func SignUpRoute(router *gin.Engine)  {
    router.POST("/signup", controllers.CreateUser()) //add this
}

func IndexRoute(router *gin.Engine, port string){
    router.LoadHTMLGlob("pkg/pages/*")
	router.GET("/api", controllers.DisplayIndexPage())
	// router.Run(":" + port)
}