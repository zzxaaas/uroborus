package model

type BaseImage struct {
	Model
	Name  string `json:"name" form:"name"`
	Tags  string `json:"tags"`
	Label string `json:"label" form:"label"`
	Port  string `json:"port" form:"port"`
}

type BaseImageReq struct {
	BaseImage
}

type GetBaseImageResp struct {
	BaseImage
	Tags []string `json:"tags" form:"tags"`
}

type RegistryImageInfoResp struct {
	Tags []struct {
		Name  string `json:"name" form:"name"`
		Layer string `json:"layer" form:"layer"`
	}
}
