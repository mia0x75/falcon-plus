package model

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | metric      | varchar(128)     | NO   |     |         |                |
// | tags        | varchar(256)     | NO   |     |         |                |
// | max_step    | int(11)          | NO   |     | 1       |                |
// | priority    | tinyint(4)       | NO   |     | 0       |                |
// | func        | varchar(16)      | NO   |     | all(#1) |                |
// | op          | varchar(8)       | NO   |     |         |                |
// | right_value | varchar(64)      | NO   |     | NULL    |                |
// | note        | varchar(128)     | NO   |     |         |                |
// | run_begin   | varchar(16)      | NO   |     |         |                |
// | run_end     | varchar(16)      | NO   |     |         |                |
// | template_id | int(10) unsigned | NO   | MUL | 0       |                |
// +-------------+------------------+------+-----+---------+----------------+

// Strategy 模型定义和映射
type Strategy struct {
	ID         int64  `json:"id"         gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	Metric     string `json:"metric"     gorm:"column:metric;size:100;type:varchar(100);not null"`
	Tags       string `json:"tags"       gorm:"column:tags;size:200;type:varchar(200);not null"`
	MaxStep    int    `json:"maxStep"    gorm:"column:max_step;type:int;not null;default:1"`
	Priority   int    `json:"priority"   gorm:"column:priority;type:tinyint;not null;default:0"`
	Func       string `json:"func"       gorm:"column:func;size:20;type:varchar(20);not null;default:'all(#1)'"`
	Op         string `json:"op"         gorm:"column:op;size:8;type:varchar(8);not null"`
	RightValue string `json:"rightValue" gorm:"column:right_value;size:50;type:varchar(50);not null"`
	Note       string `json:"note"       gorm:"column:note;size:200;type:varchar(200);not null"`
	RunBegin   string `json:"runBegin"   gorm:"column:run_begin;size:16;type:varchar(16);not null"`
	RunEnd     string `json:"runEnd"     gorm:"column:run_end;size:16;type:varchar(16);not null"`
	TemplateID int64  `json:"templateID" gorm:"column:template_id;type:int;not null;default:0;index:index_1"`
}

// TableName 结构体映射到的物理表名称
func (m Strategy) TableName() string {
	return "strategies"
}
