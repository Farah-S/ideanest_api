package handlers

import (
	"html/template"
	"net/http"

	// "github.com/example/golang-test/pkg/controllers"
	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/gin-gonic/gin"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("signup.html"))
	tmpl.Execute(w, nil)
}


func SignupSubmitHandler() gin.HandlerFunc{
    return func(c *gin.Context) {
		var user models.OrganizationMember
		if err := c.ShouldBind(&user); err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
			return
		}

		// You can perform validation here if needed

		// Save user to database or perform any other action

		// Redirect to a success page
		c.Redirect(http.StatusSeeOther, "/success")
	}
}
