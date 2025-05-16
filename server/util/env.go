package util

import "os"

// IsRelease 是否是发布环境
var IsRelease = os.Getenv("GIN_MODE") == "release"
