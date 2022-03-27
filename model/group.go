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
	UserCount    int    `json:"user_count" form:"user_count"`
	ProjectCount int    `json:"project_count" form:"project_count"`
}

type GetGroupProjectResp struct {
	Group
	Projects []GetProjectResp `json:"projects" form:"projects"`
}

type RegisterGroupProjectReq struct {
	ProjectId uint `json:"project_id" form:"project_id"`
	GroupId   uint `json:"group_id" form:"group_id"`
}
