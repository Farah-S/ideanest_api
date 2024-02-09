package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "user-athentication-golang/database"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/example/golang-test/pkg/database/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails
type SignedInDetails struct {
    Email      string
    Name 	string
    ID        primitive.ObjectID
	AccessLevel string
    jwt.StandardClaims
}

var memberCollection *mongo.Collection = repository.GetCollection(repository.DB, "organization_members", "api_db")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateAllTokens(email string, name string, id primitive.ObjectID, accessLevel string) (signedToken string, signedRefreshToken string, err error) {
    claims := &SignedInDetails{
        Email:      email,
        Name: 		name,
        ID:        	id,
		AccessLevel: accessLevel,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
        },
    }

    refreshClaims := &SignedInDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
        },
    }

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

    if err != nil {
        log.Panic(err)
        return
    }

    return token, refreshToken, err
}

//ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *SignedInDetails, msg string) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedInDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )

    if err != nil {
        msg = err.Error()
        return
    }

    claims, ok := token.Claims.(*SignedInDetails)
    if !ok {
        msg = fmt.Sprintf("the token is invalid")
        msg = err.Error()
        return
    }

    if claims.ExpiresAt < time.Now().Local().Unix() {
        msg = fmt.Sprintf("token is expired")
        msg = err.Error()
        return
    }

    return claims, msg
}

//UpdateAllTokens renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId primitive.ObjectID) {
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

    var updateObj primitive.D

    updateObj = append(updateObj, bson.E{"token", signedToken})
    updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

    Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

    upsert := true
    filter := bson.M{"_id": userId}
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    _, err := memberCollection.UpdateOne(
        ctx,
        filter,
        bson.D{
            {"$set", updateObj},
        },
        &opt,
    )
    defer cancel()

    if err != nil {
        log.Panic(err)
        return
    }

    return
}

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

	if token.ExpiresAt.Before(time.Now()) {
		return nil, mongo.ErrNoDocuments
	}

	return &token, nil
}