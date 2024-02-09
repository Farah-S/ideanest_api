package controllers

import (
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
	Id       primitive.ObjectID `bson:"_id" json:"organization_id,omitempty"`
	Name                string               `form:"name" json:"name,omitempty"  validate:"required"`
	Description         string               `form:"description" json:"description,omitempty"  validate:"required"`
	OrganizationMembers []MemberResponse `json:"organization_members,omitempty"`
}

type MemberResponse struct {
	Name                string               `json:"name,omitempty"`
	Email         string               `json:"email,omitempty"`
	AccessLevel string `bson:"access_level,omitempty" json:"access_level,omitempty"`
}

