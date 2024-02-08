package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine)  {
    //All routes related to users comes here
}

func IndexRoute(server *gin.Engine, port string){
    server.LoadHTMLGlob("pkg/api/pages/*")
	server.GET("/api", func(c *gin.Context) {
	
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
		},
	)
	// c.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})
	// apirouter := server.Group("/api")
	// apirouter.GET("/healthchecker", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	// })
	// router.Run()
	server.Run(":" + port)
}