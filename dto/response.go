package dto

type Response struct {
	Error       string `json:"error,omitempty"`
	Valid       bool   `json:"valid,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
}
