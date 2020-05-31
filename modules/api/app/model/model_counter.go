package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-------------+------------------+------+-----+------------------+----------------+
// | Field       | Type             | Null | Key | Default          | Extra          |
// +-------------+------------------+------+-----+------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | endpoint_id | int(10) unsigned | NO   | MUL | NULL             |                |
// | counter     | varchar(200)     | NO   |     |                  |                |
// | step        | int(11)          | NO   |     | 60               |                |
// | type        | varchar(16)      | NO   |     | NULL             |                |
// | ts          | int(11)          | YES  |     | NULL             |                |
// | create_at   | int(10) unsigned | NO   |     | unix_timestamp() |                |
// | update_at   | int(10) unsigned | YES  |     | NULL             |                |
// +-------------+------------------+------+-----+------------------+----------------+

// Counter 模型定义和映射
type Counter struct {
	ID         uint   `json:"id"         gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	EndpointID int    `json:"endpointID" gorm:"column:endpoint_id;type:int;not null"`
	Counter    string `json:"counter"    gorm:"column:counter;size:200;type:varchar(200);not null"`
	Step       int    `json:"step"       gorm:"column:step;type:int;not null;default:60"`
	Type       string `json:"type"       gorm:"column:type;size:16;type:varchar(16);not null"`
	Ts         int    `json:"ts"         gorm:"column:ts;type:int"`
	CreateAt   int64  `json:"createAt"   gorm:"column:create_at;type:timestamp;not null"`
	UpdateAt   int64  `json:"updateAt"   gorm:"column:update_at;type:timestamp"`
}

// TableName 结构体映射到的物理表名称
func (m Counter) TableName() string {
	return "counters"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Counter) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Counter) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
