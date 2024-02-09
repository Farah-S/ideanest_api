package routes

import (
	"net/http"

	"github.com/example/golang-test/pkg/api/handlers"
	"github.com/example/golang-test/pkg/api/middleware"
	"github.com/example/golang-test/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func StartApp(router *gin.Engine)  {
	SignUpRoute(router)
	SignInRoute(router)
	OrganizationRoutes(router)
	router.Run(":8080")
}

func OrganizationRoutes(router *gin.Engine)  {
	router.GET("/api/create-org", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_org.html", nil)
	})

	router.POST("/organization/:organization_id/invite", middleware.AuthMiddleware(),controllers.InviteUser())
	router.POST("/organization", middleware.AuthMiddleware(),controllers.CreateOrg())
	router.PUT("/organization/:organization_id", middleware.AuthMiddleware(),controllers.UpdateOrganization())
	router.DELETE("/organization/:organization_id", middleware.AuthMiddleware(),controllers.DeleteOrganization())
	router.GET("/organization/:organization_id", middleware.AuthMiddleware(),controllers.GetOrganization())
	router.GET("/organization", middleware.AuthMiddleware(),controllers.GetAllOrganizations())

}

func SignUpRoute(router *gin.Engine)  {
    router.GET("/api/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	}) 
	router.POST("/signup", controllers.CreateUser()) 
	
}

func SignInRoute(router *gin.Engine)  {
	 router.GET("/api/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", nil)
	}) 
	router.POST("/signin", controllers.GetUser()) 
	router.POST("/refresh-token", handlers.RefreshTokenHandler())
}

func IndexRoute(router *gin.Engine, port string){
    router.LoadHTMLGlob("pkg/pages/*")
	router.GET("/api", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			},
		)
	})

}