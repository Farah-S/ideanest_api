package controllers

import (
	"context"
	"encoding/json"
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
	// "go.mongodb.org/mongo-driver/mongo/options"
)
var tokenCollection *mongo.Collection = repository.GetCollection(repository.DB, "tokens", "api_db")
var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")
var validate = validator.New()

// SignedDetails
type SignedInUser struct {
    Email      string
    Name 	string
    ID        primitive.ObjectID
	AccessLevel string
	BearerToken string
	Invites		[]primitive.ObjectID
}

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
		
		filter := bson.M{"email": user.Email}

		// Perform a find operation to check if the email exists
		var existingUser models.OrganizationMember
		err := memberCollection.FindOne(ctx, filter).Decode(&existingUser)
	
		if err != mongo.ErrNoDocuments {
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
			Invites: []primitive.ObjectID{},
        }
		
        _, err = memberCollection.InsertOne(ctx, newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError,  MessageResponse{Message: "insert error"})
            return
        }

        refreshToken:= utils.GenerateRefreshToken(user.Id)

		_, err = tokenCollection.InsertOne(ctx, refreshToken)
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
        var foundToken utils.Token

        if err := c.ShouldBind(&user); err != nil {
            // c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            c.JSON(http.StatusBadRequest, TokensResponse{Message: "error "+err.Error(), AccessToken: "", RefreshToken: ""})
        	defer cancel()
			return
        }

        err := memberCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
        defer cancel()
        if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, TokensResponse{Message: "incorrect email "+err.Error(), AccessToken: "", RefreshToken: ""})
			return
        }

        passwordIsValid := utils.CheckPasswordHash(user.Password, foundUser.Password)
        defer cancel()
        if passwordIsValid != true {
            // c.JSON(http.StatusInternalServerError, "incorrect password")
            c.JSON(http.StatusBadRequest, TokensResponse{Message: "incorrect password "+err.Error(), AccessToken: "", RefreshToken: ""})
			return
        }

        // token, refreshToken, _ := utils.GenerateAllTokens(*&foundUser.Email, *&foundUser.Name, foundUser.Id, foundUser.AccessLevel)
		err = tokenCollection.FindOne(ctx, bson.M{"member_id": foundUser.Id}).Decode(&foundToken)
        defer cancel()
        if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, TokensResponse{Message: "refresh token error "+err.Error(), AccessToken: "", RefreshToken: foundToken.Token})
			return
        }
		token:=utils.GenerateAccessToken(foundUser.Id)
		
        // refreshToken:=utils.UpdateRefreshToken(foundToken.ID, foundUser.Id)
		// Serialize the user object to JSON
		cookieUser:=SignedInUser{
			ID: foundUser.Id,
			Name: foundUser.Name,
			Email: foundUser.Email,
			BearerToken: token,
			AccessLevel: foundUser.AccessLevel,
			Invites: foundUser.Invites,
		}
		userJSON, err := json.Marshal(cookieUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize user data"})
			return
		}

		// Save the serialized user data in a cookie
		c.SetCookie("user", string(userJSON), 3600, "/", "", false, false)

		c.JSON(http.StatusOK, TokensResponse{Message: "Success", AccessToken: token, RefreshToken: foundToken.Token})

        // c.JSON(http.StatusOK, foundUser)
		// c.Redirect(http.StatusOK, "/api")
    }
}