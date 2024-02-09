package controllers

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageResponse struct {
	Message string `json:"message"`
}

type TokensResponse struct {
	Message      string `json:"message" form:"message"`
	AccessToken  string `json:"access_token" form:"access_token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type IDResponse struct {
	OrganizationID primitive.ObjectID `json:"organization_id"`
}