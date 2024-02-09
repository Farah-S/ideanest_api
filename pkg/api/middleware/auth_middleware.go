package middleware

import (
	"fmt"
	"net/http"

	// helper "user-athentication-golang/helpers"

	"github.com/example/golang-test/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientToken := c.Request.Header.Get("token")
        if clientToken == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
            c.Abort()
            return
        }

        claims, err := utils.ValidateToken(clientToken)
        if err != "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            c.Abort()
            return
        }

        c.Set("email", claims.Email)
        c.Set("name", claims.Name)
        c.Set("id", claims.ID)
        c.Set("access_level", claims.AccessLevel)

        c.Next()

    }
}
