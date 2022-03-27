package model

const (
	RepoBasePath       = "/repo/"
	DockerfileBasePath = "/Dockerfile/"
	LogfileBasePath    = "/log/"
	ImgfileBasePath    = "/img/"
)

const (
	RedisKeyLogPrefix = "log"
)

type GetProjectResp struct {
	Project       `json:"project"`
	DeployHistory `json:"deploy_history"`
	Group         Group `json:"group"`
}

type Project struct {
	Model
	Image       string `json:"image" form:"image"`
	Port        string `json:"port" form:"port"`
	BindPort    int    `json:"bind_port" form:"bind_port"`
	RemoteRepo  string `json:"remote_repo" form:"remote_repo"`
	AccessUrl   string `json:"access_url" form:"access_url"`
	Branch      string `json:"branch" form:"branch"`
	LocalRepo   string `json:"local_repo" form:"local_repo"`
	Name        string `json:"name" form:"name"`
	Version     string `json:"version" form:"version"`
	Command     string `json:"command" form:"command"`
	Env         string `json:"env"`
	Container   string `json:"container" form:"container"`
	UserName    string `json:"user_name"`
	Dockerfile  string `json:"dockerfile" form:"dockerfile"`
	Type        string `json:"type" form:"type"`
	GroupId     uint   `json:"group_id" form:"group_id"`
	IsShow      int    `json:"is_show" form:"is_show"`
	SurfacePath string `json:"surface_path" form:"surface_path"`
}

type RegisterProjectReq struct {
	Project
	Env []string `json:"env" form:"env"`
	Cmd []string `json:"cmd" form:"cmd"`
}
