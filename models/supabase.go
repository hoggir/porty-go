package models

type SupRpcErrorResponse struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}
