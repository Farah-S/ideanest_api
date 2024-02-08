// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    organization, err := UnmarshalOrganization(bytes)
//    bytes, err = organization.Marshal()

package main

import "encoding/json"

func UnmarshalOrganization(data []byte) (Organization, error) {
	var r Organization
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Organization) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Organization struct {
	ID                  ID                   `json:"_id"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	OrganizationMembers []OrganizationMember `json:"organization_members"`
}

