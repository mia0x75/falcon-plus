package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +------------+------------------+------+-----+-------------------+-----------------------------+
// | Field      | Type             | Null | Key | Default           | Extra                       |
// +------------+------------------+------+-----+-------------------+-----------------------------+
// | id         | int(11) unsigned | NO   | PRI | NULL              | auto_increment              |
// | pid        | int(11) unsigned | NO   | MUL | 0                 |                             |
// | name       | char(128)        | NO   |     | NULL              |                             |
// | create_at  | timestamp        | NO   |     | CURRENT_TIMESTAMP |                             |
// | update_at  | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +------------+------------------+------+-----+-------------------+-----------------------------+

// Screen 模型定义和映射
type Screen struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	PID      int64  `json:"PID"      gorm:"column:pid;type:int;not null;index:index_1;unique:unique_1"`
	Name     string `json:"name"     gorm:"column:name;size:200;type:varchar(200);not null;unique:unique_1"`
	UpdateAt int64  `json:"updateAt" gorm:"column:update_at;type:timestamp"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Screen) TableName() string {
	return "screens"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Screen) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Screen) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
