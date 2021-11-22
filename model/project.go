package model

type Project struct {
	Model
	Language    string `json:"language" form:"language"`
	LangVersion string `json:"lang_version" form:"lang_version"`
	Port        int    `json:"port" form:"port"`
	BindPort    int    `json:"bind_port" form:"bind_port"`
	RemoteRepo  string `json:"remote_repo" form:"remote_repo"`
	Branch      string `json:"branch" form:"branch"`
	LocalRepo   string `json:"local_repo" form:"local_repo"`
	Name        string `json:"name" form:"name"`
	Version     string `json:"version" form:"version"`
	Command     string `json:"command" form:"command"`
	Container   string `json:"container" form:"container"`
	Status      string `json:"status" form:"status"`
	UserName    string `json:"user_name"`
}
