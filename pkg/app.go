package main

import (
	"net/http"

	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/gin-gonic/gin"
)

func app() {
	r := gin.Default()
	r.GET("/api/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	r.POST("/signup", func(c *gin.Context) {
		var user models.OrganizationMember
		if err := c.ShouldBind(&user); err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
			return
		}

		// You can perform validation here if needed

		// Save user to database or perform any other action

		// Redirect to a success page
		c.Redirect(http.StatusSeeOther, "/success")
	})

	r.GET("/success", func(c *gin.Context) {
		c.HTML(http.StatusOK, "success.html", nil)
	})

	r.Run(":8080")
}