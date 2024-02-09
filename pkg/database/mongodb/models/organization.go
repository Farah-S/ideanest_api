// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    organization, err := UnmarshalOrganization(bytes)
//    bytes, err = organization.Marshal()

package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalOrganization(data []byte) (Organization, error) {
	var r Organization
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Organization) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Organization struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Name                string               `form:"name" json:"name,omitempty"  validate:"required"`
	Description         string               `form:"description" json:"description,omitempty"  validate:"required"`
	OrganizationMembersIDs []primitive.ObjectID `bson:"organization_members" json:"organization_members,omitempty"`
}

