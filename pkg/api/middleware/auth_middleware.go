package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/example/golang-test/pkg/controllers"
	"github.com/gin-gonic/gin"
)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the request cookie
		userJSON, err := c.Cookie("user")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized please signin"})
			c.Abort()
			return
		}

		// Validate the token (replace with your token validation logic)
		var user controllers.SignedInUser
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user data in the request context for later use
		c.Set("user", user)

		// Continue with the next middleware or route handler
		c.Next()
	
	}
}
