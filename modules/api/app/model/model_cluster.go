package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +-------------+------------------+------+-----+-------------------+-----------------------------+
// | Field       | Type             | Null | Key | Default           | Extra                       |
// +-------------+------------------+------+-----+-------------------+-----------------------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment              |
// | group_id    | int(11)          | NO   |     | NULL              |                             |
// | numerator   | varchar(10240)   | NO   |     | NULL              |                             |
// | denominator | varchar(10240)   | NO   |     | NULL              |                             |
// | endpoint    | varchar(255)     | NO   |     | NULL              |                             |
// | metric      | varchar(255)     | NO   |     | NULL              |                             |
// | tags        | varchar(255)     | NO   |     | NULL              |                             |
// | ds_type     | varchar(255)     | NO   |     | NULL              |                             |
// | step        | int(11)          | NO   |     | NULL              |                             |
// | update_at   | timestamp        | NO   |     |                   |                             |
// | creator     | varchar(255)     | NO   |     | NULL              |                             |
// | create_at   | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +-------------+------------------+------+-----+-------------------+-----------------------------+

// Cluster 模型定义和映射
type Cluster struct {
	ID          int64  `json:"id"          gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	GroupID     int64  `json:"groupID"     gorm:"column:group_id;type:int;not null"`
	Numerator   string `json:"numerator"   gorm:"column:numerator;type:text"`
	Denominator string `json:"denominator" gorm:"column:denominator;type:text"`
	Endpoint    string `json:"endpoint"    gorm:"column:endpoint;size:200;type:varchar(200)"`
	Metric      string `json:"metric"      gorm:"column:metric;size:200;type:varchar(200)"`
	Tags        string `json:"tags"        gorm:"column:tags;size:200;type:varchar(200)"`
	DsType      string `json:"ds_type"     gorm:"column:ds_type;size:200;type:varchar(200)"`
	Step        int    `json:"step"        gorm:"column:step;type:int;not null"`
	UpdateAt    int64  `json:"updateAt"    gorm:"column:update_at"`
	Creator     int64  `json:"creator"     gorm:"column:creator;type:int;not null"`
	CreateAt    int64  `json:"createAt"    gorm:"column:create_at;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Cluster) TableName() string {
	return "clusters"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Cluster) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 插入数据前保证数据的完整性
func (m *Cluster) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
