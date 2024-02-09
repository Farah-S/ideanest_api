package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

// SignedDetails
type SignedInUser struct {
    Email      string
    Name 	string
    ID        primitive.ObjectID
	AccessLevel string
	BearerToken string
	Invites		[]primitive.ObjectID
}