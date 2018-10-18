package http

import (
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
)

func SetupIndexRoutes() {
	// 触发索引全量更新, 同步操作
	routes.GET("/index/updateAll", func(c *gin.Context) {
		go index.UpdateIndexAllByDefaultStep()
		cutils.JSONR(c, 200, gin.H{"msg": "ok"})
	})

	// 获取索引全量更新的并行数
	routes.GET("/index/updateAll/concurrent", func(c *gin.Context) {
		cutils.JSONR(c, 200, gin.H{"msg": "ok", "value": index.GetConcurrentOfUpdateIndexAll()})
	})
}
