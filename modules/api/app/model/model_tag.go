package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-------------+------------------+------+-----+------------------+----------------+
// | Field       | Type             | Null | Key | Default          | Extra          |
// +-------------+------------------+------+-----+------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | tag         | varchar(200)     | NO   | MUL |                  |                |
// | endpoint_id | int(10) unsigned | NO   |     | NULL             |                |
// | ts          | int(11)          | YES  |     | NULL             |                |
// | create_at   | int(10) unsigned | NO   |     | unix_timestamp() |                |
// | update_at   | int(10) unsigned | YES  |     | NULL             |                |
// +-------------+------------------+------+-----+------------------+----------------+

// Tag 模型定义和映射
type Tag struct {
	ID         uint   `json:"id"         gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Tag        string `json:"tag"        gorm:"column:tag;size:200;type:varchar(200);not null;unique:unique_1"`
	EndpointID int    `json:"endpointID" gorm:"column:endpoint_id;type:int;unique:unique_1"`
	Ts         int    `json:"ts"         gorm:"column:ts;type:int"`
	CreateAt   int64  `json:"createAt"   gorm:"column:create_at;type:timestamp;not null"`
	UpdateAt   int64  `json:"updateAt"   gorm:"column:update_at;type:timestamp"`
}

// TableName 结构体映射到的物理表名称
func (m Tag) TableName() string {
	return "tags"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Tag) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 更新数据前保证数据的完整性
func (m *Tag) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
