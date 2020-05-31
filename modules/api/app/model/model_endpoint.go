package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-----------+------------------+------+-----+------------------+----------------+
// | Field     | Type             | Null | Key | Default          | Extra          |
// +-----------+------------------+------+-----+------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | endpoint  | varchar(200)     | NO   | UNI |                  |                |
// | ts        | int(11)          | YES  |     | NULL             |                |
// | create_at | int(10) unsigned | NO   |     | unix_timestamp() |                |
// | update_at | int(10) unsigned | YES  |     | NULL             |                |
// +-----------+------------------+------+-----+------------------+----------------+

// Endpoint 模型定义和映射
type Endpoint struct {
	ID       uint      `json:"id"       gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Endpoint string    `json:"endpoint" gorm:"column:endpoint;size:200;type:varchar(200);not null;unique:unique_1"`
	Ts       int       `json:"-"        gorm:"column:ts;type:int"`
	CreateAt int64     `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
	UpdateAt int64     `json:"updateAt" gorm:"column:update_at;type:timestamp"`
	Counters []Counter `json:"-"        gorm:"ForeignKey:EndpointIDE"`
}

// TableName 结构体映射到的物理表名称
func (m Endpoint) TableName() string {
	return "endpoints"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Endpoint) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Endpoint) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
