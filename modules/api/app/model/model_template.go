package model

import (
	"time"
)

// +-----------+------------------+------+-----+------------------+----------------+
// | Field     | Type             | Null | Key | Default          | Extra          |
// +-----------+------------------+------+-----+------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | name      | varchar(75)      | NO   | UNI |                  |                |
// | parent_id | int(10) unsigned | NO   |     | 0                |                |
// | action_id | int(10) unsigned | NO   |     | 0                |                |
// | creator   | int(10) unsigned | NO   | MUL | NULL             |                |
// | create_at | int(10) unsigned | NO   |     | unix_timestamp() |                |
// +-----------+------------------+------+-----+------------------+----------------+

// Template 模型定义和映射
type Template struct {
	ID       int64  `json:"id"        gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	Name     string `json:"name"      gorm:"column:name;size:75;type:varchar(75);not null"`
	ParentID int64  `json:"parent_id" gorm:"column:parent_id;type:int;not null;default:0"`
	ActionID int64  `json:"action_id" gorm:"column:action_id;type:int;not null;default:0"`
	Creator  int64  `json:"creator"   gorm:"column:creator;type:int;not null"`
	CreateAt int64  `json:"createAt"  gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Template) TableName() string {
	return "templates"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m Template) BeforeCreate() (err error) {
	m.CreateAt = time.Now().Unix()
	return
}
