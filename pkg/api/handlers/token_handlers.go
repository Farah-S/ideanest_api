package handlers

import (
	"net/http"

	"github.com/example/golang-test/pkg/utils"
	"github.com/gin-gonic/gin"
)

func refreshTokenHandler(c *gin.Context) {
	// Parse the refresh token from the request
	refreshToken := c.PostForm("refresh_token")

	// Check if the refresh token is valid
	token, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate a new refresh token
	newToken := utils.GenerateRefreshToken(token.MemberID)

	// Store the new refresh token in the database
	utils.InsertToken(newToken)

	c.JSON(http.StatusOK, gin.H{
		"refresh_token": newToken.Token,
	})
}
