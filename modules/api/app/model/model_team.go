package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+------------------+------+-----+------------------+----------------+
// | Field     | Type             | Null | Key | Default          | Extra          |
// +-----------+------------------+------+-----+------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | name      | varchar(50)      | NO   | UNI | NULL             |                |
// | resume    | varchar(200)     | NO   |     |                  |                |
// | creator   | int(10) unsigned | NO   |     | 0                |                |
// | create_at | int(10) unsigned | NO   |     | unix_timestamp() |                |
// +-----------+------------------+------+-----+------------------+----------------+

// Team 模型定义和映射
type Team struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;auto_increment;primary_key;not null"`
	Name     string `json:"name"     gorm:"column:name;unique_index:unique_1;size:50;type:varchar(50);not null"`
	Resume   string `json:"resume"   gorm:"column:resume;size:200;type:varchar(200);not null"`
	Creator  int64  `json:"creator"  gorm:"column:creator;type:int;not null"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Team) TableName() string {
	return "teams"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Team) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
