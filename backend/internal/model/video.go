package model

type VideoGenRequest struct {
	Model         string           `json:"model"`
	ModelOverride string           `json:"model_override"`
	Prompt        string           `json:"prompt"`
	Duration      int              `json:"duration"`
	Resolution    string           `json:"resolution"`
	Ratio         string           `json:"ratio"`
	Seed          *int             `json:"seed"`
	Watermark     *bool            `json:"watermark"`
	FirstFrame    string           `json:"first_frame"`
	RefImages     []string         `json:"ref_images"`
	VideoURL      string           `json:"video_url"`
	AudioSetting  string           `json:"audio_setting"`
	Media         []map[string]any `json:"media"`
	Input         map[string]any   `json:"input"`
	Parameters    map[string]any   `json:"parameters"`
}
