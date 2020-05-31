package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+------------------+------+-----+-------------------+----------------+
// | Field     | Type             | Null | Key | Default           | Extra          |
// +-----------+------------------+------+-----+-------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | name      | varchar(50)      | NO   | UNI |                   |                |
// | creator   | varchar(64)      | NO   |     |                   |                |
// | create_at | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// | come_from | tinyint(4)       | NO   |     | 0                 |                |
// +-----------+------------------+------+-----+-------------------+----------------+

// Group 模型定义和映射
type Group struct {
	ID       int64  `json:"id"        gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Name     string `json:"name"      gorm:"column:name;size:75;type:varchar(75);not null;unique:unique_1"`
	Creator  int64  `json:"creator"   gorm:"column:creator;type:int;not null"`
	ComeFrom int    `json:"-"         gorm:"column:come_from;type:tinyint;not null;default:0"`
	CreateAt int64  `json:"create_at" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Group) TableName() string {
	return "groups"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Group) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	scope.SetColumn("ComeFrom", 1) // TODO: ComeFrom ??
	return
}
