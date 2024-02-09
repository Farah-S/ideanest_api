// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    organizationMember, err := UnmarshalOrganizationMember(bytes)
//    bytes, err = organizationMember.Marshal()

package models

import "encoding/json"
import "go.mongodb.org/mongo-driver/bson/primitive"

func UnmarshalOrganizationMember(data []byte) (OrganizationMember, error) {
	var r OrganizationMember
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *OrganizationMember) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type OrganizationMember struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name        string `json:"name,omitempty" validate:"required"`
	Email       string `json:"email,omitempty" validate:"required"`
	Password    string `json:"password,omitempty" validate:"required"`
	AccessLevel string `json:"access_level"`
}

