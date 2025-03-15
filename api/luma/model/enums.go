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

type ImageModelType string

const (
	ModelPhoton1     ImageModelType = "photon-1"
	ModelPhotonFlash ImageModelType = "photon-flash-1"
)

type ImageReference struct {
	URL    string  `json:"url"`
	Weight float64 `json:"weight"`
}

type StyleReference struct {
	URL    string  `json:"url"`
	Weight float64 `json:"weight"`
}

type CharacterRef struct {
	Identity0 struct {
		Images []string `json:"images"`
	} `json:"identity0"`
}

type ModifyImageRef struct {
	URL    string  `json:"url"`
	Weight float64 `json:"weight"`
}

var (
	ValidAspectRatios = []string{
		"1:1", "3:4", "4:3",
		"9:16", "16:9", "9:21", "21:9",
	}
	ValidImageModels = []ImageModelType{
		ModelPhoton1, ModelPhotonFlash,
	}
)
