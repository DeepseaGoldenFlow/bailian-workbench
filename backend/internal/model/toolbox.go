package model

type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type OCRRequest struct {
	ImageBase64 string `json:"image_base64"`
	ImageURL    string `json:"image_url"`
}

type CodeRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type DocumentRequest struct {
	Task     string `json:"task"`
	Text     string `json:"text"`
	Question string `json:"question"`
}