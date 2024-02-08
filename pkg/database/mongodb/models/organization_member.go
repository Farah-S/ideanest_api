// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    organizationMember, err := UnmarshalOrganizationMember(bytes)
//    bytes, err = organizationMember.Marshal()

package main

import "encoding/json"

func UnmarshalOrganizationMember(data []byte) (OrganizationMember, error) {
	var r OrganizationMember
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *OrganizationMember) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type OrganizationMember struct {
	ID          ID     `json:"_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessLevel string `json:"access_level"`
}

type ID struct {
	OID string `json:"$oid"`
}
