package utils

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")

var SECRET_KEY string = os.Getenv("SECRET_KEY")


func UpdateRefreshToken(refreshTokenId  primitive.ObjectID, userId primitive.ObjectID) string{
    var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	newtoken:= GenerateRefreshToken(userId)
	newtoken.ID=refreshTokenId
	// Example filter and update values
	filter := bson.M{"_id": refreshTokenId}
	update := bson.M{"$set": newtoken}

	// Perform the update operation
	_, err := tokenCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Panic(err)
	}

    return newtoken.Token
}

// Token represents a refresh token stored in the database
type Token struct {
	ID        primitive.ObjectID    `bson:"_id,omitempty"`
	Token     string    `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
	MemberID    primitive.ObjectID    `bson:"member_id"`
}

func GenerateAccessToken(memberID primitive.ObjectID) string{
	claims := jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(2)).Unix(),
        }
  

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
     if err != nil {
        log.Panic(err) 
    }
	return token
}

var tokenCollection *mongo.Collection = repository.GetCollection(repository.DB, "tokens", "api_db")

func GenerateRefreshToken(memberID primitive.ObjectID) *Token {
	refreshClaims :=  jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	}
    

    // token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
        log.Panic(err)
        
    }
	// token := "your_token_generation_logic_here"
	expiresAt :=  time.Now().Local().Add(time.Hour * time.Duration(24)) // Example: Token expires in 24 hours

	return &Token{
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		MemberID:    memberID,
	}
}

func InsertToken(token *Token) {
	_, err := tokenCollection.InsertOne(context.Background(), token)
	if err != nil {
		log.Println("Error inserting token:", err)
	}
}

func ValidateRefreshToken(tokenStr string) (*Token, error) {
	var token Token
	err := tokenCollection.FindOne(context.Background(), bson.M{"token": tokenStr}).Decode(&token)
	if err != nil {
		return nil, err
	}

	// if token.ExpiresAt.Before(time.Now()) {
	// 	return nil, mongo.ErrEmptySlice
	// }

	return &token, nil
}