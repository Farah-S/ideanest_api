package controllers


import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DisplayIndexPage() gin.HandlerFunc{
	return func(c *gin.Context) {
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
		return
	}
}


func DisplaySignUpPage() gin.HandlerFunc{
	return func(c *gin.Context) {
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"signup.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "SignUp",
			},
		)
		
		// c.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
		return
	}
}