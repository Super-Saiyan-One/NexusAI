package model

type GenerationResponse struct {
	ID             string         `json:"id"`
	GenerationType string         `json:"generation_type"`
	State          string         `json:"state"`
	FailureReason  *string        `json:"failure_reason"`
	CreatedAt      string         `json:"created_at"`
	Assets         *AssetsDetails `json:"assets"`
	Model          ModelType      `json:"model"`
	Request        RequestInfo    `json:"request"`
}

type AssetsDetails struct {
	Video         *string `json:"video,omitempty"`
	Image         *string `json:"image,omitempty"`
	ProgressVideo *string `json:"progress_video,omitempty"`
}

type RequestInfo struct {
	GenerationType string      `json:"generation_type"`
	Prompt         string      `json:"prompt"`
	AspectRatio    string      `json:"aspect_ratio"`
	Loop           bool        `json:"loop"`
	Keyframes      interface{} `json:"keyframes"`
	CallbackURL    *string     `json:"callback_url"`
	Model          ModelType   `json:"model"`
	Resolution     *Resolution `json:"resolution"`
	Duration       *Duration   `json:"duration"`
}
