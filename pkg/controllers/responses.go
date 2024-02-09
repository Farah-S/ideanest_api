package controllers

import (
	"github.com/example/golang-test/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type OneOrgResponse struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Name                string               `form:"name" json:"name,omitempty"  validate:"required"`
	Description         string               `form:"description" json:"description,omitempty"  validate:"required"`
	OrganizationMembers []models.OrganizationMember `json:"organization_members,omitempty"`
}

