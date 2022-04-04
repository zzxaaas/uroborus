package model

const (
	DontShowProjStatus = 0
	ShowProjStatus     = 1
	RedisKeyPrefix     = "group"
	KeyUserCountSuffix = "-uc"
	KeyProjCountSuffix = "-pc"
)

type Group struct {
	Model
	Name         string `json:"name" form:"name"`
	Description  string `json:"description" form:"description"`
	OriginId     uint   `json:"origin_id" form:"origin_id"`
	CreateUser   string `json:"create_user" form:"create_user"`
	UserCount    int64  `json:"user_count" form:"user_count"`
	ProjectCount int64  `json:"project_count" form:"project_count"`
	Code         string `json:"code" form:"code"`
}

type GetGroupProjectResp struct {
	Group
	Projects []GetProjectResp `json:"projects" form:"projects"`
}

type RegisterGroupProjectReq struct {
	Name string `json:"name" form:"name"`
	Code string `json:"code" form:"code"`
}
