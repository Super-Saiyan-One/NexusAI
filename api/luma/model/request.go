package model

type VideoGenerateRequest struct {
	GenerationType string      `json:"generation_type"`
	Prompt         string      `json:"prompt"`
	Model          ModelType   `json:"model"`
	AspectRatio    string      `json:"aspect_ratio,omitempty"`
	Resolution     *Resolution `json:"resolution,omitempty"`
	Duration       *Duration   `json:"duration,omitempty"`
	Loop           bool        `json:"loop,omitempty"`
	CallbackURL    *string     `json:"callback_url,omitempty"`
}
