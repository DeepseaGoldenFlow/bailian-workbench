package model

type ImageGenRequest struct {
	Model          string         `json:"model"`
	ModelOverride  string         `json:"model_override"`
	Prompt         string         `json:"prompt"`
	NegativePrompt string         `json:"negative_prompt"`
	Size           string         `json:"size"`
	N              int            `json:"n"`
	Seed           *int           `json:"seed"`
	RefImg         string         `json:"ref_img"`
	RefStrength    float64        `json:"ref_strength"`
	Steps          int            `json:"steps"`
	PromptExtend   *bool          `json:"prompt_extend"`
	Images         []string       `json:"images"`
	Input          map[string]any `json:"input"`
	Parameters     map[string]any `json:"parameters"`
}
