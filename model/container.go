package model

type BuildImageOption struct {
	Path       string
	Tag        string
	Version    string
	Dockerfile string
}

type ContainerOption struct {
	Name  string
	Image string
	// 绑定
	Port string
	// 容器内
	ProtoPort string
	Env       []string
}
