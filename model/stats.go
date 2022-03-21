package model

type ContainerStats struct {
	Read        string      `json:"read"`
	Preread     string      `json:"preread"`
	PidStats    PidsStats   `json:"pids_stats"`
	NumProcs    int         `json:"num_procs"`
	CPUStats    CPUStats    `json:"cpu_stats"`
	PrecpuStats CPUStats    `json:"precpu_stats"`
	MemoryStats MemoryStats `json:"memory_stats"`
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Networks    Networks    `json:"networks"`
}
type Networks struct {
	Eth0 Eth0 `json:"eth0"`
}
type Eth0 struct {
	RxBytes   int `json:"rx_bytes"`
	RxPackets int `json:"rx_packets"`
	RxErrors  int `json:"rx_errors"`
	RxDropped int `json:"rx_dropped"`
	TxBytes   int `json:"tx_bytes"`
	TxPackets int `json:"tx_packets"`
	TxErrors  int `json:"tx_errors"`
	TxDropped int `json:"tx_dropped"`
}

type MemoryStats struct {
	Usage    int   `json:"usage"`
	MaxUsage int   `json:"max_usage"`
	Limit    int64 `json:"limit"`
}

type CPUStats struct {
	CPUUsage       CPUUsage `json:"cpu_usage"`
	SystemCPUUsage int64    `json:"system_cpu_usage"`
	OnlineCpus     int      `json:"online_cpus"`
}

type CPUUsage struct {
	TotalUsage  int64   `json:"total_usage"`
	PercpuUsage []int64 `json:"percpu_usage"`
}
type PidsStats struct {
	Current int `json:"current"`
}
