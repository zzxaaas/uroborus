package model

type BuildImageOption struct {
	Path       string
	Tag        string
	Version    string
	Dockerfile string
	DeployID   uint
}

type ContainerOption struct {
	Name  string
	Image string
	// 绑定
	Port string
	// 容器内
	ProtoPort string
	Env       []string
	NeedPull  bool
	DeployID  uint
}

type LogStream struct {
	Stream string `json:"stream"`
}
