package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | Field        | Type             | Null | Key | Default           | Extra                       |
// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | id           | int              | NO   | PRI | NULL              | auto_increment              |
// | case_id      | varchar(50)      | YES  | MUL | NULL              |                             |
// | step         | int unsigned     | YES  |     | NULL              |                             |
// | cond         | varchar(200)     | NO   |     | NULL              |                             |
// | status       | int unsigned     | YES  |     | 0                 |                             |
// | create_at    | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +--------------+------------------+------+-----+-------------------+-----------------------------+

// Event 模型定义和映射
type Event struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	CaseID   string `json:"caseID"   gorm:"column:case_id;size:50;type:varchar(50);index:index_1"`
	Step     int    `json:"step"     gorm:"column:step;type:int"`
	Cond     string `json:"cond"     gorm:"column:cond;size:200;type:varchar(200)"`
	Status   int    `json:"status"   gorm:"column:status;type:int;not null;default:0"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Event) TableName() string {
	return "events"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Event) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
