package dto

type WebResponse struct {
	Status string `json:"status"`
	Data any `json:"data,omitempty"`
	Errs any `json:"errors,omitempty"`
}