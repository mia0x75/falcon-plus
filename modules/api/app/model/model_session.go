package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+------------------+------+-----+------------------+----------------+
// | Field     | Type             | Null | Key | Default          | Extra          |
// +-----------+------------------+------+-----+------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | user_id   | int(10) unsigned | NO   | MUL | NULL             |                |
// | sign      | char(32)         | NO   | MUL | NULL             |                |
// | expire    | int(10) unsigned | NO   |     | NULL             |                |
// | create_at | int(10) unsigned | NO   |     | unix_timestamp() |                |
// +-----------+------------------+------+-----+------------------+----------------+

// Session 模型定义和映射
type Session struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	UserID   int64  `json:"userID"   gorm:"column:user_id;unique_index:unique_1;not null"`
	Sign     string `json:"sign"     gorm:"column:sign;type:char(32);index:index_1;size:32;not null"`
	Expire   int    `json:"expire"   gorm:"column:expire;type:int;not null"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Session) TableName() string {
	return "sessions"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Session) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
