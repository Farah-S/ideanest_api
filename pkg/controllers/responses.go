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
	Message	string	`json:"message" form:"message"`
	AccessToken	string	`json:"access_token" form:"access_token"`
	RefreshToken	string	`json:"refresh_token" form:"refresh_token"`
}

type IDResponse struct{
	OrganizationID string                 `json:"_id"`
}