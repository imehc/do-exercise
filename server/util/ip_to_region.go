package util

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
)

func init() {
	dbFile := "./ip2region.xdb"
	ipBuff, err := xdb.LoadContentFromFile(dbFile)
	if err != nil {
		fmt.Printf("加载数据库数据失败 `%s`: %s\n", dbFile, err)
		return
	}
	searcher, err := xdb.NewWithBuffer(ipBuff)
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
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		global.Log.Error("解析IP地址失败: %s", zap.Error(err))
		return "-"
	}
	return region
}
