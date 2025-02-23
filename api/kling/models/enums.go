package models

// 模型版本枚举
const (
	ModelKlingV1   = "kling-v1"
	ModelKlingV1_6 = "kling-v1-6"
)

// 视频模式枚举
const (
	VideoModeStandard = "std"
	VideoModePro      = "pro"
)

// 画面比例枚举
const (
	AspectRatio16_9 = "16:9"
	AspectRatio9_16 = "9:16"
	AspectRatio1_1  = "1:1"
)

// 运镜类型枚举
const (
	CameraTypeSimple           = "simple"
	CameraTypeDownBack         = "down_back"
	CameraTypeForwardUp        = "forward_up"
	CameraTypeRightTurnForward = "right_turn_forward"
	CameraTypeLeftTurnForward  = "left_turn_forward"
)
