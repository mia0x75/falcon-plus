package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

// APIIndexItemInput TODO:
type APIIndexItemInput struct {
	Endpoint string `json:"endpoint" form:"endpoint" binding:"required"`
	Metric   string `json:"metric"   form:"metric"   binding:"required"`
	Step     int    `json:"step"     form:"step"     binding:"required"`
	Dstype   string `json:"dstype"   form:"dstype"   binding:"required"`
	Tags     string `json:"tags"     form:"tags"`
}

// SetupCommonRoutes 设置路由
func SetupCommonRoutes() {
	routes.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok\n")
	})

	routes.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("%s\n", g.Version))
	})

	routes.GET("/workdir", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("%s\n", file.SelfDir()))
	})

	routes.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, cu.Dto{Msg: "success", Data: g.Config()})
	})

	routes.GET("/config/reload", func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RemoteAddr, "127.0.0.1") {
			g.ParseConfig(g.ConfigFile)
			c.JSON(http.StatusOK, cu.Dto{Msg: "success", Data: "ok"})
		} else {
			c.JSON(http.StatusOK, cu.Dto{Msg: "success", Data: "no privilege"})
		}
	})

	routes.GET("/stats/graph-queue-size", func(c *gin.Context) {
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
		cu.JSONR(c, 200, rt)
	})

	routes.GET("/counter/migrate.v2", func(c *gin.Context) {
		counter := rrdtool.GetCounterV2()
		log.Debugf("[D] migrating counter v2: %v", counter)
		c.JSON(200, counter)
	})

	// 更新一条索引数据,用于手动建立索引 endpoint metric step dstype tags
	routes.POST("/index", func(c *gin.Context) {
		inputs := []*APIIndexItemInput{}
		if err := c.Bind(&inputs); err != nil {
			c.AbortWithError(500, err)
			return
		}

		for _, in := range inputs {
			err, tags := cu.SplitTagsString(in.Tags)
			if err != nil {
				log.Errorf("[E] split tags: %s error: %v", in.Tags, err)
				continue
			}

			err = index.UpdateIndexOne(in.Endpoint, in.Metric, tags, in.Dstype, in.Step)
			if err != nil {
				log.Errorf("[E] build index fail, item: %v, error: %v", in, err)
			} else {
				log.Debugf("[D] build index manually, item: %v", in)
			}
		}
		cu.JSONR(c, 200, gin.H{"msg": "ok"})
	})

	// Compatible with open-falcon v0.1
	routes.GET("/counter/migrate", func(c *gin.Context) {
		cnt := rrdtool.GetCounter()
		log.Debugf("[D] migrating counter: %s", cnt)
		c.JSON(http.StatusOK, gin.H{"msg": "ok", "counter": cnt})
	})
}
