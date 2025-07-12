package util

import (
	"runtime"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	B  = 1
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

type SystemInfo struct {
	Os    Os     `json:"os"`
	Cpu   Cpu    `json:"cpu"`
	Ram   Ram    `json:"ram"`
	Disks []Disk `json:"disks"`
}

type Os struct {
	GOOS         string `json:"goos"`
	NumCPU       int    `json:"num_cpu"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	Used        int `json:"used"`
	Total       int `json:"total"`
	UsedPercent int `json:"used_percent"`
}

type Disk struct {
	MountPoint  string `json:"mount_point"`
	Used        int    `json:"used"`
	Total       int    `json:"total"`
	UsedPercent int    `json:"use_dpercent"`
}

// getSystemInfo 获取系统信息
func GetSystemInfo() (s *SystemInfo, err error) {
	s = &SystemInfo{}
	s.Os.GOOS = runtime.GOOS
	s.Os.NumCPU = runtime.NumCPU()
	s.Os.Compiler = runtime.Compiler
	s.Os.GoVersion = runtime.Version()
	s.Os.NumGoroutine = runtime.NumGoroutine()

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	s.Ram.Total = int(v.Total) / MB
	s.Ram.Used = int(v.Used) / MB
	s.Ram.UsedPercent = int(v.UsedPercent)

	cores, err := cpu.Counts(false)
	if err != nil {
		return nil, err
	}
	s.Cpu.Cores = cores

	cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true)
	if err != nil {
		return nil, err
	}
	s.Cpu.Cpus = cpus

	for i := range global.Config.DiskList {
		mp := global.Config.DiskList[i].MountPoint
		u, err := disk.Usage(mp)
		if err != nil {
			return nil, err
		}
		s.Disks = append(s.Disks, Disk{
			MountPoint:  mp,
			Used:        int(u.Used) / MB,
			Total:       int(u.Total) / MB,
			UsedPercent: int(u.UsedPercent),
		})
	}

	return s, nil
}
