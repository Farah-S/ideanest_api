package controllers

import (
	"context"
	"fmt"
	// "log"
	"net/http"
	"time"

	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"

	// "github.com/example/golang-test/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var orgCollection *mongo.Collection = repository.GetCollection(repository.DB, "organizations", "api_db")

func CreateOrg() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var org models.Organization
		var currentUser, bl = c.Get("user")
		
		defer cancel()
		if !bl{
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "please log in"})
			return
		}
		if err := c.ShouldBind(&org); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{Message: "bind error " + err.Error()})
			return
		}
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

		res, err := orgCollection.InsertOne(ctx, newOrg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, MessageResponse{Message: "insert error"})
			return
		}

		c.JSON(http.StatusCreated, IDResponse{OrganizationID: res.InsertedID.(primitive.ObjectID)})
		c.Redirect(http.StatusSeeOther, "/api")
	}
}

func GetOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var id=c.Param("organization_id")
		fmt.Println(id)
		var org models.Organization
		//pass these options to the Find method
		objectID, err := primitive.ObjectIDFromHex(id)
		defer cancel()
		if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
			return
        }
		filter := bson.M{"_id": objectID}
		err = orgCollection.FindOne(ctx, filter).Decode(&org)

		// orgs, err :=orgCollection.Find(ctx,nil)
        //  memberCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
        if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
			return
        }
		
		response:=OneOrgResponse{Id: org.Id,
			Name: org.Name, 
			Description: org.Description, 
			OrganizationMembers:getMembersOfOrg(org.OrganizationMembersIDs,c)}
		
		c.JSON(http.StatusOK, response)
	}
}

func GetAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
       
		//pass these options to the Find method
		findOptions := options.Find()
		
		//Define an array in which you can store the decoded documents
		// var results []Member

		//Passing the bson.D{{}} as the filter matches  documents in the collection
    	orgs, err := orgCollection.Find(context.TODO(), bson.D{{}}, findOptions)

		// orgs, err :=orgCollection.Find(ctx,nil)
        //  memberCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
        defer cancel()
        if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
			return
        }
		var results []models.Organization
		if err := orgs.All(context.Background(), &results); err != nil {
			// panic(err)
	        c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})

		}
		response:=[]OneOrgResponse{}
		for i := 0; i < len(results); i++ {
			fmt.Println(results[i].Id)
			members:=getMembersOfOrg(results[i].OrganizationMembersIDs,c)
			response = append(response,  
				OneOrgResponse{Id: results[i].Id, 
					Name:  results[i].Name, 
					Description: results[i].Description, 
					OrganizationMembers:members })
		}
		c.JSON(http.StatusOK, response)

	}
}

func getMembersOfOrg(ids []primitive.ObjectID, c *gin.Context)  []MemberResponse{
	// Define a slice to store documents
	// var members []models.OrganizationMember
	var results []MemberResponse

	// Define the filter to find documents with IDs in the list
	filter := bson.M{"_id": bson.M{"$in": ids}}

	// Find documents in the collection that match the filter
	cursor, err := memberCollection.Find(context.Background(), filter)
	if err != nil {
        c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})

		// panic(err)
	}

	// Iterate over the cursor and decode each document into the slice
	if err := cursor.All(context.Background(), &results); err != nil {
		// panic(err)
        c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
	}

	// for i := 0; i < len(members); i++ {
	// 	results = append(results, MemberResponse{Name: members[i].Name, Email: members[i].Email, AccessLevel: members[i].AccessLevel})
	// }
	return results

}