package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+------------------+------+-----+-------------------+----------------+
// | Field     | Type             | Null | Key | Default           | Extra          |
// +-----------+------------------+------+-----+-------------------+----------------+
// | id        | int(11) unsigned | NO   | PRI | NULL              | auto_increment |
// | endpoints | text             | NO   |     |                   |                |
// | counters  | text             | NO   |     |                   |                |
// | ck        | char(32)         | NO   | UNI | NULL              |                |
// | create_at | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// +-----------+------------------+------+-----+-------------------+----------------+
// time_ => create_at
// ck => checksum char(32) or binary(16)
// Draft 模型定义和映射
type Draft struct {
	ID        int64  `json:"id"        gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Endpoints string `json:"endpoints" gorm:"column:endpoints;type:text;not null"`
	Counters  string `json:"counters"  gorm:"column:counters;type:text;not null"`
	Sign      string `json:"sign"      gorm:"column:sign;size:32;type:char(32);not null;unique:unique_1"`
	CreateAt  int64  `json:"createAt"  gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Draft) TableName() string {
	return "drafts"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Draft) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
