package http

import (
	"github.com/gin-gonic/gin"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
)

// SetupIndexRoutes 设置路由
func SetupIndexRoutes() {
	// 触发索引全量更新, 同步操作
	routes.GET("/index/updateAll", func(c *gin.Context) {
		go index.UpdateIndexAllByDefaultStep()
		cu.JSONR(c, 200, gin.H{"msg": "ok"})
	})

	// 获取索引全量更新的并行数
	routes.GET("/index/updateAll/concurrent", func(c *gin.Context) {
		cu.JSONR(c, 200, gin.H{"msg": "ok", "value": index.GetConcurrentOfUpdateIndexAll()})
	})
}
