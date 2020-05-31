package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-------------+------------------+------+-----+-------------------+----------------+
// | Field       | Type             | Null | Key | Default           | Extra          |
// +-------------+------------------+------+-----+-------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | group_id    | int(10) unsigned | NO   | MUL | NULL              |                |
// | dir         | varchar(255)     | NO   |     | NULL              |                |
// | create_user | varchar(64)      | NO   |     |                   |                |
// | create_at   | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// +-------------+------------------+------+-----+-------------------+----------------+

// Plugin 模型定义和映射
type Plugin struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	GroupID  int64  `json:"groupID"  gorm:"column:group_id;type:int;not null;index:index_1"`
	Dir      string `json:"dir"      gorm:"column:dir;size:200;type:varchar(200);not null"`
	Creator  int64  `json:"creator"  gorm:"column:creator;type:int;not null"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Plugin) TableName() string {
	return "plugin_dir"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Plugin) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
