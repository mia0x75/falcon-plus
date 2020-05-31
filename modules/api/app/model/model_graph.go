package model

import (
	"time"
)

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(11) unsigned | NO   | PRI | NULL    | auto_increment |
// | title       | char(128)        | NO   |     | NULL    |                |
// | hosts       | varchar(10240)   | NO   |     |         |                |
// | counters    | varchar(1024)    | NO   |     |         |                |
// | screen_id   | int(11) unsigned | NO   | MUL | NULL    |                |
// | timespan    | int(11) unsigned | NO   |     | 3600    |                |
// | type        | char(2)          | NO   |     | h       |                |
// | method      | char(8)          | YES  |     |         |                |
// | position    | int(11) unsigned | NO   |     | 0       |                |
// | tags        | varchar(512)     | NO   |     |         |                |
// +-------------+------------------+------+-----+---------+----------------+

// Graph 模型定义和映射
type Graph struct {
	ID       int64  `json:"id"       gorm:"column:id;type:int;not null;auto_increment;primary_key"`
	Title    string `json:"title"    gorm:"column:title;size:100;type:varchar(100);not null"`
	Hosts    string `json:"hosts"    gorm:"column:hosts;type:text;not null"`
	Counters string `json:"counters" gorm:"column:counters;type:text;not null"`
	ScreenID int64  `json:"screenID" gorm:"column:screen_id;type:int;not null;index:index_1"`
	TimeSpan int    `json:"timespan" gorm:"column:timespan;type:int;not null;default:3600"`
	Type     string `json:"type"     gorm:"column:type;size:2;type:char(2);not null; default:'h'"`
	Method   string `json:"method"   gorm:"column:method;size:8;type:varchar(8)"`
	Position int    `json:"position" gorm:"column:position;type:int;not null;default:0"`
	Tags     string `json:"tags"     gorm:"column:tags;size:500;type:varchar(500);not null"`
	CreateAt int64  `json:"createAt" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m Graph) TableName() string {
	return "graphs"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m Graph) BeforeCreate() (err error) {
	m.CreateAt = time.Now().Unix()
	return
}
