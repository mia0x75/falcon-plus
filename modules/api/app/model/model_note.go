package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | Field        | Type             | Null | Key | Default           | Extra                       |
// +--------------+------------------+------+-----+-------------------+-----------------------------+
// | id           | mediumint(9)     | NO   | PRI | NULL              | auto_increment              |
// | event_caseId | varchar(50)      | YES  | MUL | NULL              |                             |
// | note         | varchar(300)     | YES  |     | NULL              |                             |
// | case_id      | varchar(20)      | YES  |     | NULL              |                             |
// | status       | varchar(15)      | YES  |     | NULL              |                             |
// | timestamp    | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// | user_id      | int(10) unsigned | YES  | MUL | NULL              |                             |
// +--------------+------------------+------+-----+-------------------+-----------------------------+

// Note 模型定义和映射
type Note struct {
	ID          int64  `json:"id"          gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	EventCaseID string `json:"eventCaseID" gorm:"column:event_caseId;size:50;type:varchar(50);index:index_1"`
	Note        string `json:"note"        gorm:"column:note;size:200;type:varchar(200)"`
	CaseID      string `json:"caseID"      gorm:"column:case_id;size:20;type:varchar(20)"`
	Status      string `json:"status"      gorm:"column:status;size:15;type:varchar(15)"`
	CreateAt    int64  `json:"createAt"    gorm:"column:create_at;type:timestamp;not null"`
	Creator     int64  `json:"creator"     gorm:"column:creator;type:int"`
}

// TableName 结构体映射到的物理表名称
func (m Note) TableName() string {
	return "notes"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Note) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}
