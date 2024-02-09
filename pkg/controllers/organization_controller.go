package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var orgCollection *mongo.Collection = repository.GetCollection(repository.DB, "organizations", "api_db")

func CreateOrg() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var org models.Organization
		var currentUser, bl = c.Get("user")
		if !bl{
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "please log in"})
			return
		}
		if err := c.ShouldBind(&org); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "bind error " + err.Error()})
			return
		}
		defer cancel()
		// log.Fatal(c.)
		//validate the request body
		if err := c.ShouldBind(&org); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "bind error " + err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&org); validationErr != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "validator error"})
			return
		}

		filter := bson.M{"name": org.Name}

		// Perform a find operation to check if the email exists
		var existingOrg models.OrganizationMember
		err := orgCollection.FindOne(ctx, filter).Decode(&existingOrg)

		if err != mongo.ErrNoDocuments {
			// c.Redirect(http.StatusSeeOther,"/api/signup")
			c.JSON(http.StatusInternalServerError, MessageResponse{Message: "organization already exists"})
			return
		}
		user:=currentUser.(SignedInUser)
		ids:=[]primitive.ObjectID{user.ID}
		newOrg := models.Organization{
			Id:          primitive.NewObjectID(),
			Name:        org.Name,
			OrganizationMembersIDs:      ids,
			Description: org.Description,
		}

		_, err = orgCollection.InsertOne(ctx, newOrg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, MessageResponse{Message: "insert error"})
			return
		}

		c.JSON(http.StatusCreated, IDResponse{OrganizationID: "success"})
		c.Redirect(http.StatusSeeOther, "/api")
	}
}