package model

type VideoGenerationRequest struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt,omitempty"`
	PromptOptimizer  *bool    `json:"prompt_optimizer,omitempty"`
	FirstFrameImage  string   `json:"first_frame_image,omitempty"`
	SubjectReference []string `json:"subject_reference,omitempty"`
	CallbackURL      string   `json:"callback_url,omitempty"`
}

type TaskResponse struct {
	TaskID   string   `json:"task_id"`
	BaseResp BaseResp `json:"base_resp"`
}

type TaskStatus struct {
	TaskID      string   `json:"task_id"`
	Status      string   `json:"status"`
	FileID      string   `json:"file_id,omitempty"`
	VideoWidth  int      `json:"video_width,omitempty"`
	VideoHeight int      `json:"video_height,omitempty"`
	BaseResp    BaseResp `json:"base_resp"`
}

type FileResponse struct {
	FileURL    string   `json:"file_url"`
	ExpireTime int64    `json:"expire_time"`
	BaseResp   BaseResp `json:"base_resp"`
}

type BaseResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CallbackRequest struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	FileID    string `json:"file_id,omitempty"`
	Challenge string `json:"challenge"`
}

type CallbackChallenge struct {
	Challenge string `json:"challenge"`
}
