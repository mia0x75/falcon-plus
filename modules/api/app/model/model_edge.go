package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/open-falcon/falcon-plus/modules/api/g"
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
// Edge 模型定义和映射
// 1: Team -> User
// 2: Group -> Host
// 3: Group -> Template
type Edge struct {
	ID           int64 `json:"id"           gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Type         int   `json:"type"         gorm:"column:type;type:tinyint;not null"`
	AncestorID   int64 `json:"ancestorID"   gorm:"column:ancestor_id;type:text;not null"`
	DescendantID int64 `json:"descendantID" gorm:"column:descendant_id;type:int;not null"`
	CreateAt     int64 `json:"createAt"     gorm:"column:create_at;type:timestamp;not null"`
	Creator      int64 `json:"creator"      gorm:"column:creator;type:int;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Edge) TableName() string {
	return "edges"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Edge) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// Exists 判断主机群组关联是否存在
func (m Edge) Exists() bool {
	var r *Edge
	db := g.Con()
	if db.Where("group_id = ? AND host_id = ? AND type = ?", m.AncestorID, m.DescendantID, m.Type).Find(r); r != nil {
		return true
	} else {
		return false
	}
}
