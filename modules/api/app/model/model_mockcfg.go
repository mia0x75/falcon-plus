package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+---------------------+------+-----+-------------------+-----------------------------+
// | Field     | Type                | Null | Key | Default           | Extra                       |
// +-----------+---------------------+------+-----+-------------------+-----------------------------+
// | id        | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment              |
// | name      | varchar(255)        | NO   | UNI |                   |                             |
// | obj       | varchar(10240)      | NO   |     |                   |                             |
// | obj_type  | varchar(255)        | NO   |     |                   |                             |
// | metric    | varchar(128)        | NO   |     |                   |                             |
// | tags      | varchar(1024)       | NO   |     |                   |                             |
// | dstype    | varchar(32)         | NO   |     | GAUGE             |                             |
// | step      | int(11) unsigned    | NO   |     | 60                |                             |
// | mock      | double              | NO   |     | 0                 |                             |
// | creator   | varchar(64)         | NO   |     |                   |                             |
// | update_at | int                 | NO   |     | NULL              |                             |
// | create_at | int                 | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +-----------+---------------------+------+-----+-------------------+-----------------------------+

// Mockcfg 模型定义和映射
type Mockcfg struct {
	ID       int64   `json:"id"       gorm:"column:id;type:int;auto_increment;primary_key;not null"`
	Name     string  `json:"name"     gorm:"column:name;size:200;type:varchar(200);not null;unique:unique_1"`
	Obj      string  `json:"obj"      gorm:"column:obj;type:text;not null"`
	ObjType  string  `json:"obj_type" gorm:"column:obj_type;size:200;type:varchar(200);not null"` // group, host, other
	Metric   string  `json:"metric"   gorm:"column:metric;size:128;type:varchar(128);not null"`
	Tags     string  `json:"tags"     gorm:"column:tags;size:500;type:varchar(500);not null"`
	DsType   string  `json:"dstype"   gorm:"column:dstype;size:32;type:varchar(32);not null"`
	Step     int     `json:"step"     gorm:"column:step;type:int"`
	Mock     float64 `json:"mock"     gorm:"column:mock"`
	Creator  int64   `json:"creator"  gorm:"column:creator;type:int;not null"`
	CreateAt int64   `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
	UpdateAt int64   `json:"updateAt" gorm:"column:update_at;type:timestamp"`
}

// TableName 结构体映射到的物理表名称
func (m Mockcfg) TableName() string {
	return "mockcfg"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Mockcfg) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Mockcfg) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
