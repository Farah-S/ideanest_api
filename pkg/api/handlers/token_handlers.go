package handlers

import (
	"net/http"

	"github.com/example/golang-test/pkg/controllers"
	"github.com/example/golang-test/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RefreshTokenHandler() gin.HandlerFunc{
    return func(c *gin.Context) {
		
		refresh_token:=c.Query("refresh_token")
		
		token, err := utils.ValidateRefreshToken(refresh_token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"refresh_token": err.Error()})
			return
		}

		// Generate a new refresh token
		newToken := utils.UpdateRefreshToken(token.ID,token.MemberID)
		accessToken:=utils.GenerateAccessToken(token.MemberID)

		c.JSON(http.StatusOK, controllers.TokensResponse{AccessToken: accessToken, RefreshToken: newToken,Message: "Success"})
	}
}
