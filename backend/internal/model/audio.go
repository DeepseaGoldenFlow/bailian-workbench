package model

type TTSRequest struct {
	Input  string `json:"input"`
	Voice  string `json:"voice"`
	Format string `json:"format"`
}