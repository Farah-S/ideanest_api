package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/example/golang-test/pkg/database/mongodb/models"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"github.com/example/golang-test/pkg/utils"
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
		user:=currentUser.(utils.SignedInUser)
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
	
	// var members []models.OrganizationMember
	var results []MemberResponse

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := memberCollection.Find(context.Background(), filter)
	if err != nil {
        c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})

		// panic(err)
	}

	if err := cursor.All(context.Background(), &results); err != nil {
		// panic(err)
        c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
	}

	// for i := 0; i < len(members); i++ {
	// 	results = append(results, MemberResponse{Name: members[i].Name, Email: members[i].Email, AccessLevel: members[i].AccessLevel})
	// }
	return results
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonid:=c.Param("organization_id") 
		id, err:=primitive.ObjectIDFromHex(jsonid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if(!utils.IsMember(id,c)){
			c.JSON(http.StatusUnauthorized, "You are not a member thus can't update")
			return
		}
		name:=c.Query("name")
		description:=c.Query("description") 
		// Define filter to match the document to update
		filter := bson.M{"_id": id}

		// Define update document to specify the modifications
		update := bson.M{
			"$set": bson.M{
				"name": name,
				"description": description,
			},
		}

		// Perform the update operation
		res, err := orgCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if(res.ModifiedCount==0){
			c.JSON(http.StatusBadRequest, "error updating")
			return
		}
		c.JSON(http.StatusOK, OneOrgResponse{Id: id, Name: name,Description: description})

	}
}


func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonid:=c.Param("organization_id") 
		id, err:=primitive.ObjectIDFromHex(jsonid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if(!utils.IsMember(id,c)){
			c.JSON(http.StatusUnauthorized, "You are not a member thus can't delete")
			return
		}
		// Define filter to match the document to update
		filter := bson.M{"_id": id}
		
		// Perform the update operation
		_, err = orgCollection.DeleteOne(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, MessageResponse{Message: "Success"})

	}
}

func InviteUser()  gin.HandlerFunc{
    return func(c *gin.Context) {
		userEmail:=c.Query("user_email")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var id=c.Param("organization_id")
		
		objectID, err := primitive.ObjectIDFromHex(id)
		defer cancel()
		if err != nil {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect email"})
            c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
			return
        }
		filter := bson.M{"email": userEmail}
		var user models.OrganizationMember
		err = memberCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
            c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
			return
        }
		user.Invites = append(user.Invites, objectID)
		UpdateMember(user,c)
        c.JSON(http.StatusOK, MessageResponse{Message: "Success"})

	}
}