package index

import (
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

// 初始化索引功能模块
func Start() {
	cfg := g.Config()
	if !cfg.Index.Enabled {
		log.Info("[I] index.Start warning, not enable")
		return
	}

	InitDB()
	if cfg.Index.AutoDelete {
		StartIndexDeleteTask()
		log.Info("[I] index.Start warning, index cleaner enable")
	}
	go StartIndexUpdateAllTask()
	log.Info("[I] index.Start ok")
}
