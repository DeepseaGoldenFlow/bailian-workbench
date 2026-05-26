package model

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}