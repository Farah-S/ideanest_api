package controllers

type UserResponse struct {
    Status  int                    `json:"status"`
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}

type MessageResponse struct{
	Message string                 `json:"message"`
}

type TokensResponse struct{
	Message	string
	AccessToken	string
	RefreshToken	string
}

type IDResponse struct{
	OrganizationID string                 `json:"_id"`
}