package model

const (
	RedisKeyCtnPrefix  = "container"
	KeyCtnPreCpuSuffix = "-cpu"
	KeyCtnPreSysSuffix = "-sys"
)

type ConnectContainerReq struct {
	ID string `json:"id" form:"id"`
}

type GetContainerStatsReq struct {
	ID           string  `json:"id" form:"id"`
	PreCPUUsage  float64 `json:"pre_cpu_usage" form:"pre_cpu_usage"`
	PresSysUsage float64 `json:"pres_sys_usage" form:"pres_sys_usage"`
}

type GetContainerStatsResp struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	CPUUsage   float64 `json:"cpu_usage"`
	MemUsage   uint64  `json:"mem_usage"`
	MemLimit   uint64  `json:"mem_limit"`
	MemPercent float64 `json:"mem_percent"`
	NETIn      uint64  `json:"net_in"`
	NETOut     uint64  `json:"net_out"`
}

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
	Status string `json:"status"`
	Id     string `json:"id"`
}
