package server

import "time"

type Service struct {
	Name        string    `json:"name,omitempty"`
	Ip          string    `json:"ip,omitempty"`
	RefreshTime time.Time `json:"refresh_time"`
	Info        Info      `json:"info"`
}

type Info struct {
	MemAll          uint64  `json:"mem_all,omitempty"`
	MemFree         uint64  `json:"mem_free,omitempty"`
	MemUsed         uint64  `json:"mem_used,omitempty"`
	MemUsedPercent  float64 `json:"mem_used_percent,omitempty"`
	CpuPercent      float64 `json:"cpu_percent,omitempty"`
	DiskTotal       uint64  `json:"disk_total"`
	DiskUsed        uint64  `json:"disk_used"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
	IOBytesSent     uint64  `json:"io_bytes_sent"`
	IOBytesRev      uint64  `json:"io_bytes_rev"`
}
