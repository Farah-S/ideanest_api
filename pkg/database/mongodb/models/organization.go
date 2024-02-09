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
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name                string               `json:"name,omitempty"  validate:"required"`
	Description         string               `json:"description,omitempty"  validate:"required"`
	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty"  validate:"required"`
}

