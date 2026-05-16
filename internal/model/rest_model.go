package model

type RestRequest struct {
	Endpoint string `json:"endpoint,omitzero"`
	Method   string `json:"method,omitzero"`
}

type RestResponse struct {
	Endpoint   string `json:"endpoint,omitzero"`
	Method     string `json:"method,omitzero"`
	StatusCode string `json:"statusCode,omitzero"`
}
