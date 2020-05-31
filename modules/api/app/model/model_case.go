package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | Field          | Type             | Null | Key | Default           | Extra                       |
// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | id             | varchar(50)      | NO   | PRI | NULL              |                             |
// | endpoint       | varchar(100)     | NO   | MUL | NULL              |                             |
// | metric         | varchar(200)     | NO   |     | NULL              |                             |
// | func           | varchar(50)      | YES  |     | NULL              |                             |
// | cond           | varchar(200)     | NO   |     | NULL              |                             |
// | note           | varchar(500)     | YES  |     | NULL              |                             |
// | max_step       | int(10) unsigned | YES  |     | NULL              |                             |
// | current_step   | int(10) unsigned | YES  |     | NULL              |                             |
// | priority       | int(6)           | NO   |     | NULL              |                             |
// | status         | varchar(20)      | NO   |     | NULL              |                             |
// | create_at      | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// | update_at      | timestamp        | YES  |     | NULL              |                             |
// | closed_at      | timestamp        | YES  |     | NULL              |                             |
// | closed_note    | varchar(250)     | YES  |     | NULL              |                             |
// | user_modified  | int(10) unsigned | YES  |     | NULL              |                             |
// | tpl_creator    | varchar(64)      | YES  |     | NULL              |                             |
// | expression_id  | int(10) unsigned | YES  |     | NULL              |                             |
// | strategy_id    | int(10) unsigned | YES  |     | NULL              |                             |
// | template_id    | int(10) unsigned | YES  |     | NULL              |                             |
// | process_note   | mediumint(9)     | YES  |     | NULL              |                             |
// | process_status | varchar(20)      | YES  |     | unresolved        |                             |
// +----------------+------------------+------+-----+-------------------+-----------------------------+

// Case 模型定义和映射
type Case struct {
	ID            string `json:"id"            gorm:"column:id;size:50;type:varchar(50);not null;primary_key"`
	Endpoint      string `json:"endpoint"      gorm:"column:endpoint;size:100;type:varchar(100);not null;index:index_1"`
	Metric        string `json:"metric"        gorm:"column:metric;size:200;type:varchar(200);not null"`
	Func          string `json:"func"          gorm:"column:func;size:50;type:varchar(50)"`
	Cond          string `json:"cond"          gorm:"column:cond;size:200;type:varchar(200);not null"`
	Note          string `json:"note"          gorm:"column:note;size:200;type:varchar(200)"`
	MaxStep       int    `json:"maxStep"       gorm:"column:max_step;type:int"`
	CurrentStep   int    `json:"currentStep"   gorm:"column:current_step;type:int"`
	Priority      int    `json:"priority"      gorm:"column:priority;type:int;not null"`
	Status        string `json:"status"        gorm:"column:status;size:20;type:varchar(20);not null"`
	CreateAt      int64  `json:"createAt"      gorm:"column:create_at;type:timestamp;not null"`
	UpdateAt      int64  `json:"updateAt"      gorm:"column:update_at;type:timestamp"`
	ClosedAt      int64  `json:"closedAt"      gorm:"column:closed_at;type:timestamp"`
	ClosedNote    string `json:"closedNote"    gorm:"column:closed_note;size:200;type:varchar(200)"`
	UserModified  int64  `json:"userModified"  gorm:"column:user_modified;type:int"`
	ExpressionID  int64  `json:"expressionID"  gorm:"column:expression_id;type:int"`
	StrategyID    int64  `json:"strategyID"    gorm:"column:strategy_id;type:int;index:index_1"`
	TemplateID    int64  `json:"templateID"    gorm:"column:template_id;type:int;index:index_1"`
	ProcessNote   int64  `json:"processNote"   gorm:"column:process_note;type:mediumint"`
	ProcessStatus string `json:"processStatus" gorm:"column:process_status;size:20;type:varchar(20);default:'unresolved'"`
}

// TableName 结构体映射到的物理表名称
func (m Case) TableName() string {
	return "cases"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m *Case) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return
}

// BeforeUpdate 插入数据前保证数据的完整性
func (m *Case) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return
}
