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
	Id       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Name        string `form:"name,omitempty" json:"name,omitempty" validate:"required"`
	Email       string `form:"email,omitempty" json:"email,omitempty" validate:"required"`
	Password    string `form:"password,omitempty" json:"password,omitempty" validate:"required"`
	AccessLevel string `form:"access_level" json:"access_level" bson:"access_level"`
	Invites []primitive.ObjectID `json:"invites,omitempty"`
}


