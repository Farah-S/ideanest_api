package routes

import (
	"net/http"

	"github.com/example/golang-test/pkg/api/handlers"
	"github.com/example/golang-test/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func StartApp(router *gin.Engine)  {
	// r := gin.Default()
	// r.GET("/api", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	SignUpRoute(router)
	SignInRoute(router)
	// r.GET("/success", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "success.html", nil)
	// })

	router.Run(":8080")
}

func SignUpRoute(router *gin.Engine)  {
    router.GET("/api/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	}) //add this
	router.POST("/signup", controllers.CreateUser()) //add this
	// SignUpRoute(router)
	// router.GET("/success", func(c *gin.Context) {
    //     c.HTML(http.StatusOK, "index.html", nil)
    // })
}

func SignInRoute(router *gin.Engine)  {
	 router.GET("/api/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", nil)
	}) //add this
	router.POST("/signin", controllers.GetUser()) //add this
	router.POST("/refresh-token", handlers.RefreshTokenHandler())
}

func IndexRoute(router *gin.Engine, port string){
    router.LoadHTMLGlob("pkg/pages/*")
	router.GET("/api", controllers.DisplayIndexPage())
	// router.Run(":" + port)
}