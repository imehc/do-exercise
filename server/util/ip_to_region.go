package util

import (
	"fmt"
	"os"

	"github.com/imehc/do-exercise/server/global"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
)

func init() {
	dbFile := "./ip2region.xdb"
	// 文件不存在时静默跳过，避免测试或非部署环境下刷噪音。
	if _, err := os.Stat(dbFile); err != nil {
		return
	}
	ipBuff, err := xdb.LoadContentFromFile(dbFile)
	if err != nil {
		fmt.Printf("加载数据库数据失败 `%s`: %s\n", dbFile, err)
		return
	}
	header, err := xdb.LoadHeaderFromBuff(ipBuff)
	if err != nil {
		fmt.Printf("读取数据库头信息失败 `%s`: %s\n", dbFile, err)
		return
	}
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		fmt.Printf("识别数据库IP版本失败 `%s`: %s\n", dbFile, err)
		return
	}
	searcher, err := xdb.NewWithBuffer(version, ipBuff)
	if err != nil {
		fmt.Printf("创建searcher失败: %s\n", err.Error())
		return
	}
	global.Searcher = searcher
}

// IPToRegion 解析IP地址
func IPToRegion(ip string) string {
	searcher := global.Searcher
	if searcher == nil {
		global.Log.Error("searcher初始化失败")
		return "-"
	}
	region, err := searcher.Search(ip)
	if err != nil {
		global.Log.Error("解析IP地址失败", zap.Error(err))
		return "-"
	}
	return region
}
