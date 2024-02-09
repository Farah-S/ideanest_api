package controllers

import (
	"context"
	// "log"
	"net/http"
	"time"

	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"github.com/example/golang-test/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var user models.OrganizationMember
        defer cancel()
		// log.Fatal(c.)
        //validate the request body
        if err := c.ShouldBind(&user); err != nil {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: "bind error "+err.Error()})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&user); validationErr != nil {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: "validator error"})
            return
        }
		
		// Create a unique index on the 'email' field
		_, err := memberCollection.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			// c.Redirect(http.StatusSeeOther,"/api/signup")
			c.JSON(http.StatusInternalServerError, MessageResponse{Message: "email already exists"})
			return
		}

		hashedPass := user.Password
		hashedPass, err = utils.HashPassword(hashedPass)
        
		if err != nil {
            c.JSON(http.StatusInternalServerError, MessageResponse{Message: "hash error"})
            return
        }

		newUser := models.OrganizationMember{
            Id:       primitive.NewObjectID(),
            Name:     user.Name,
            Email: user.Email,
			Password: hashedPass,
			AccessLevel: user.AccessLevel,
        }
		
        _, err = memberCollection.InsertOne(ctx, newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError,  MessageResponse{Message: "insert error"})
            return
        }
		
        c.JSON(http.StatusCreated,  MessageResponse{Message: "success"})
		c.Redirect(http.StatusSeeOther,"/api")
    }
}


func GetUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var user models.OrganizationMember
        var foundUser models.OrganizationMember

        if err := c.ShouldBind(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        err := memberCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            return
        }

        passwordIsValid := utils.CheckPasswordHash(user.Password, foundUser.Password)
        defer cancel()
        if passwordIsValid != true {
            c.JSON(http.StatusInternalServerError, "incorrect password")
            return
        }

        token, refreshToken, _ := utils.GenerateAllTokens(*&foundUser.Email, *&foundUser.Name, foundUser.Id, foundUser.AccessLevel)

        utils.UpdateAllTokens(token, refreshToken, foundUser.Id, foundUser.AccessLevel)

        // c.JSON(http.StatusOK, foundUser)
		c.Redirect(http.StatusOK, "/api")
    }
}