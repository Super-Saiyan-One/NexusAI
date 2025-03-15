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

type ImageGenerateRequest struct {
	Prompt         string           `json:"prompt"`
	AspectRatio    string           `json:"aspect_ratio,omitempty"`
	Model          ImageModelType   `json:"model,omitempty"`
	ImageRef       []ImageReference `json:"image_ref,omitempty"`
	StyleRef       []StyleReference `json:"style_ref,omitempty"`
	CharacterRef   *CharacterRef    `json:"character_ref,omitempty"`
	ModifyImageRef *ModifyImageRef  `json:"modify_image_ref,omitempty"`
}
