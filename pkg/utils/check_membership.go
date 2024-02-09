package utils

import (
	"context"
	"encoding/json"
	"net/http"

	// "github.com/example/golang-test/pkg/controllers"
	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orgCollection *mongo.Collection = repository.GetCollection(repository.DB, "organizations", "api_db")
// var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")

func IsMember(orgId primitive.ObjectID,c *gin.Context) bool {
	userJSON, err := c.Cookie("user")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please signin"})
		c.Abort()
		return false
	}

	// Validate the token (replace with your token validation logic)
	var user SignedInUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return false
	}
	
	// Define the filter to find documents with the specified value in the array field
	filter := bson.M{"_id": orgId}
	org:=models.Organization{}
	// Find documents in the collection that match the filter
	err = orgCollection.FindOne(context.Background(), filter).Decode(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error finding organization: "+err.Error())
		return false
	}
	for i := 0; i < len(org.OrganizationMembersIDs); i++ {
		if(org.OrganizationMembersIDs[i] == user.ID){
			return true
		}
	}
	return false
}