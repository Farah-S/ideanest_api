package routes

import (
	"net/http"

	"github.com/example/golang-test/pkg/api/handlers"
	"github.com/example/golang-test/pkg/api/middleware"
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
	OrganizationRoutes(router)
	// r.GET("/success", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "success.html", nil)
	// })

	router.Run(":8080")
}

func OrganizationRoutes(router *gin.Engine)  {
	router.GET("/api/create-org", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_org.html", nil)
	})
	// router.GET("/api/find-organization", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "find_org.html", nil)
	// })

	router.POST("/organization/:organization_id/invite", middleware.AuthMiddleware(),controllers.InviteUser())
	router.POST("/organization", middleware.AuthMiddleware(),controllers.CreateOrg())
	router.PUT("/organization/:organization_id", middleware.AuthMiddleware(),controllers.UpdateOrganization())
	router.GET("/organization/:organization_id", middleware.AuthMiddleware(),controllers.GetOrganization())
	router.GET("/organization", middleware.AuthMiddleware(),controllers.GetAllOrganizations())

	// router.GET("/organization", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "show_all_org.html", nil)
	// })
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