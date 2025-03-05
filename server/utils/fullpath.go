package utils

import "fmt"

// FormatFullpath 格式化全路径
func FormatFullpath(parentId, id uint, path string) string {
	if parentId == 0 {
		return fmt.Sprintf("/0/%d", id)
	}
	return fmt.Sprintf("%s/%d", path, id)
}
