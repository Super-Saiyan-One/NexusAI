package models

// TextToVideoRequest 文生视频请求参数
type TextToVideoRequest struct {
	ModelName      string         `json:"model_name,omitempty" validate:"omitempty,oneof=kling-v1 kling-v1-6"`
	Prompt         string         `json:"prompt" validate:"required,max=2500"`
	NegativePrompt string         `json:"negative_prompt,omitempty" validate:"max=2500"`
	CfgScale       float32        `json:"cfg_scale,omitempty" validate:"omitempty,min=0,max=1"`
	Mode           string         `json:"mode,omitempty" validate:"omitempty,oneof=std pro"`
	CameraControl  *CameraControl `json:"camera_control,omitempty"`
	AspectRatio    string         `json:"aspect_ratio,omitempty" validate:"omitempty,oneof=16:9 9:16 1:1"`
	Duration       string         `json:"duration,omitempty" validate:"omitempty,oneof=5 10"`
	CallbackURL    string         `json:"callback_url,omitempty" validate:"omitempty,url"`
	ExternalTaskID string         `json:"external_task_id,omitempty" validate:"omitempty,max=64"`
}

// CameraControl 摄像机控制参数
type CameraControl struct {
	Type   string        `json:"type,omitempty" validate:"omitempty,oneof=simple down_back forward_up right_turn_forward left_turn_forward"`
	Config *CameraConfig `json:"config,omitempty" validate:"required_if=Type simple"`
}

// CameraConfig 摄像机运动配置
type CameraConfig struct {
	Horizontal float32 `json:"horizontal,omitempty" validate:"excluded_with_all=Vertical Pan Tilt Roll Zoom,omitempty,min=-10,max=10"`
	Vertical   float32 `json:"vertical,omitempty" validate:"excluded_with_all=Horizontal Pan Tilt Roll Zoom,omitempty,min=-10,max=10"`
	Pan        float32 `json:"pan,omitempty" validate:"excluded_with_all=Horizontal Vertical Tilt Roll Zoom,omitempty,min=-10,max=10"`
	Tilt       float32 `json:"tilt,omitempty" validate:"excluded_with_all=Horizontal Vertical Pan Roll Zoom,omitempty,min=-10,max=10"`
	Roll       float32 `json:"roll,omitempty" validate:"excluded_with_all=Horizontal Vertical Pan Tilt Zoom,omitempty,min=-10,max=10"`
	Zoom       float32 `json:"zoom,omitempty" validate:"excluded_with_all=Horizontal Vertical Pan Tilt Roll,omitempty,min=-10,max=10"`
}

// ImageRequest文生图视频请求参数
type TextToImageRequest struct {
	Model          string  `json:"model,omitempty" validate:"omitempty,oneof=kling-v1"`
	Prompt         string  `json:"prompt" validate:"required,max=500"`
	NegativePrompt string  `json:"negative_prompt,omitempty" validate:"max=200"`
	Image          string  `json:"image,omitempty"` // Base64或URL
	ImageFidelity  float64 `json:"image_fidelity,omitempty" validate:"omitempty,min=0,max=1"`
	N              int     `json:"n,omitempty" validate:"omitempty,min=1,max=9"`
	AspectRatio    string  `json:"aspect_ratio,omitempty" validate:"omitempty,oneof=16:9 9:16 1:1 4:3 3:4 3:2 2:3"`
	CallbackURL    string  `json:"callback_url,omitempty" validate:"omitempty,url"`
}

// 验证组标签
const (
	WhenTypeSimple = "required_if=Type simple"
)
