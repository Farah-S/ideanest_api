package controllers

import (
    "context"
    "net/http"
    "time"
	
	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"github.com/example/golang-test/pkg/utils"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var user models.OrganizationMember
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: "error"})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&user); validationErr != nil {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: "error"})
            return
        }
		if user.AccessLevel==""{
			user.AccessLevel="user"
		}
		hashedPass := user.Password
		hashedPass, err := utils.HashPassword(hashedPass)
        
		if err != nil {
            c.JSON(http.StatusInternalServerError, MessageResponse{Message: "error"})
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
            c.JSON(http.StatusInternalServerError,  MessageResponse{Message: "error"})
            return
        }
		
        c.JSON(http.StatusCreated,  MessageResponse{Message: "success"})
    }
}