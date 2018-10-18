package http

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	"github.com/toolkits/file"
)

type APIIndexItemInput struct {
	Endpoint string `json:"endpoint" form:"endpoint" binding:"required"`
	Metric   string `json:"metric" form:"metric" binding:"required"`
	Step     int    `json:"step" form:"step" binding:"required"`
	Dstype   string `json:"dstype" form:"dstype" binding:"required"`
	Tags     string `json:"tags" form:"tags"`
}

func SetAPIRoutes() {
	routes.GET("/api/v2/health", func(c *gin.Context) {
		cutils.JSONR(c, 200, gin.H{"msg": "ok"})
	})

	routes.GET("/api/v2/version", func(c *gin.Context) {
		cutils.JSONR(c, 200, gin.H{"value": g.VERSION})
	})

	routes.GET("/api/v2/workdir", func(c *gin.Context) {
		cutils.JSONR(c, 200, gin.H{"value": file.SelfDir()})
	})

	routes.GET("/api/v2/config", func(c *gin.Context) {
		cutils.JSONR(c, 200, gin.H{"value": g.Config()})
	})

	routes.POST("/api/v2/config/reload", func(c *gin.Context) {
		g.ParseConfig(g.ConfigFile)
		cutils.JSONR(c, 200, gin.H{"msg": "ok"})
	})

	routes.GET("/api/v2/stats/graph-queue-size", func(c *gin.Context) {
		rt := make(map[string]int)
		for i := 0; i < store.GraphItems.Size; i++ {
			keys := store.GraphItems.KeysByIndex(i)
			oneHourAgo := time.Now().Unix() - 3600

			count := 0
			for _, ckey := range keys {
				item := store.GraphItems.First(ckey)
				if item == nil {
					continue
				}

				if item.Timestamp > oneHourAgo {
					count++
				}
			}
			i_s := strconv.Itoa(i)
			rt[i_s] = count
		}
		cutils.JSONR(c, 200, rt)
	})

	routes.GET("/api/v2/counter/migrate", func(c *gin.Context) {
		counter := rrdtool.GetCounterV2()
		log.Debug("migrating counter v2:", fmt.Sprintf("%+v", counter))
		c.JSON(200, counter)
	})

	// 更新一条索引数据,用于手动建立索引 endpoint metric step dstype tags
	routes.POST("/api/v2/index", func(c *gin.Context) {
		inputs := []*APIIndexItemInput{}
		if err := c.Bind(&inputs); err != nil {
			c.AbortWithError(500, err)
			return
		}

		for _, in := range inputs {
			err, tags := cutils.SplitTagsString(in.Tags)
			if err != nil {
				log.Error("split tags:", in.Tags, "error:", err)
				continue
			}

			err = index.UpdateIndexOne(in.Endpoint, in.Metric, tags, in.Dstype, in.Step)
			if err != nil {
				log.Error("build index fail, item:", in, "error:", err)
			} else {
				log.Debug("build index manually", in)
			}
		}
		cutils.JSONR(c, 200, gin.H{"msg": "ok"})
	})
}
