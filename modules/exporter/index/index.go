package index

import (
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

// 初始化索引功能模块
func Start() {
	cfg := g.Config()
	if !cfg.Index.Enabled {
		log.Info("[I] index.Start warning, not enable")
		return
	}

	if err := InitDB(); err != nil {
		log.Errorf("[E] open database connection fail: %v", err)
		return
	}
	defer Close()

	if cfg.Index.AutoDelete {
		StartIndexDeleteTask()
		log.Info("[I] index.Start warning, index cleaner enable")
	}
	go StartIndexUpdateAllTask()
	log.Info("[I] index.Start ok")
}
