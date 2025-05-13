package system

import "github.com/imehc/do-exercise/server/util"

type SysInfoService struct{}

// Get 获取系统信息
func (s *SysInfoService) Get() (*util.SystemInfo, error) {
	return util.GetSystemInfo()
}
