package model

type ModelType string

const (
	ModelRay2   ModelType = "ray-2"
	ModelRay1_6 ModelType = "ray-1-6"
)

type Resolution string

const (
	Res480p  Resolution = "480p"
	Res720p  Resolution = "720p"
	Res1080p Resolution = "1080p"
	Res4k    Resolution = "4k"
)

type Duration string

const (
	Duration2s Duration = "2s"
	Duration4s Duration = "4s"
	Duration5s Duration = "5s"
)

var (
	ValidModels      = []ModelType{ModelRay2, ModelRay1_6}
	ValidResolutions = []Resolution{Res480p, Res720p, Res1080p, Res4k}
	ValidDurations   = []Duration{Duration2s, Duration4s, Duration5s}
)
