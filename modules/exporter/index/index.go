// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package index

import (
	"log"

	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

// 初始化索引功能模块
func Start() {
	cfg := g.Config()
	if !cfg.Index.Enable {
		log.Println("index.Start warning, not enable")
		return
	}

	InitDB()
	if cfg.Index.AutoDelete {
		StartIndexDeleteTask()
		log.Println("index.Start warning, index cleaner enable")
	}
	StartIndexUpdateAllTask()
	log.Println("index.Start ok")
}
