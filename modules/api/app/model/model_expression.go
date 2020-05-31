package model

// +-------------+------------------+------+-----+---------+----------------+
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int unsigned     | NO   | PRI | NULL    | auto_increment |
// | expression  | varchar(1024)    | NO   |     | NULL    |                |
// | func        | varchar(16)      | NO   |     | all(#1) |                |
// | op          | varchar(8)       | NO   |     |         |                |
// | right_value | varchar(16)      | NO   |     |         |                |
// | max_step    | int              | NO   |     | 1       |                |
// | priority    | tinyint          | NO   |     | 0       |                |
// | note        | varchar(1024)    | NO   |     |         |                |
// | action_id   | int unsigned     | NO   |     | 0       |                |
// | create_user | varchar(64)      | NO   |     |         |                |
// | pause       | tinyint          | NO   |     | 0       |                |
// +-------------+------------------+------+-----+---------+----------------+

// Expression 模型定义和映射
type Expression struct {
	ID         int64  `json:"id"         gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	Expression string `json:"expression" gorm:"column:expression;size:500;type:varchar(500);not null"`
	Func       string `json:"func"       gorm:"column:func;size:50;type:varchar(50);not null;default:'all(#1)'"`
	Op         string `json:"op"         gorm:"column:op;size:8;type:varchar(8);not null"`
	RightValue string `json:"rightValue" gorm:"column:right_value;size:20;type:varchar(20);not null"`
	MaxStep    int    `json:"maxStep"    gorm:"column:max_step;type:int;not null;default:1"`
	Priority   int    `json:"priority"   gorm:"column:priority;type:tinyint;not null;default:0"`
	Note       string `json:"note"       gorm:"column:note;size:200;type:varchar(200);not null"`
	ActionID   int64  `json:"actionID"   gorm:"column:action_id;type:int;not null;default:0"`
	Creator    int64  `json:"creator"    gorm:"column:creator;type:int;not null"`
	Pause      int    `json:"pause"      gorm:"column:pause;type:tinyint;not null;default:0"`
}

// TableName 结构体映射到的物理表名称
func (m Expression) TableName() string {
	return "expressions"
}
