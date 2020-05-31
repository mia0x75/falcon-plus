package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | Field          | Type             | Null | Key | Default           | Extra                       |
// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | id             | int(11)          | NO   | PRI | NULL              | auto_increment              |
// | hostname       | varchar(255)     | NO   | UNI |                   |                             |
// | ip             | varchar(16)      | NO   |     |                   |                             |
// | agent_version  | varchar(16)      | NO   |     |                   |                             |
// | plugin_version | varchar(128)     | NO   |     |                   |                             |
// | maintain_begin | int(10) unsigned | NO   |     | 0                 |                             |
// | maintain_end   | int(10) unsigned | NO   |     | 0                 |                             |
// | update_at      | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +----------------+------------------+------+-----+-------------------+-----------------------------+

// Host 模型定义和映射
type Host struct {
	ID            int64  `json:"id"            gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Hostname      string `json:"hostname"      gorm:"column:hostname;size:200;type:varchar(200);not null;unique:unique_1"`
	IP            string `json:"IP"            gorm:"column:ip;size:15;type:varchar(15);not null"`
	AgentVersion  string `json:"agentVersion"  gorm:"column:agent_version;size:20;type:varchar(20);not null"`
	PluginVersion string `json:"pluginVersion" gorm:"column:plugin_version;size:20;type:varchar(20);not null"`
	MaintainBegin int64  `json:"maintainBegin" gorm:"column:maintain_begin;type:int;not null;default:0"`
	MaintainEnd   int64  `json:"maintainEnd"   gorm:"column:maintain_end;type:int;not null;default:0"`
	UpdateAt      int64  `json:"updateAt"      gorm:"column:update_at;type:timestamp"`
	CreateAt      int64  `json:"createAt"      gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Host) TableName() string {
	return "hosts"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Host) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Host) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
