package model

type ContainerOption struct {
	Name  string
	Image string
	// 绑定
	Port string
	// 容器内
	ProtoPort  string
	Dockerfile string
	Env        []string
}
