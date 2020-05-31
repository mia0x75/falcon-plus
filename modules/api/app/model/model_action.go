package model

// +----------------------+------------------+------+-----+---------+----------------+
// | Field                | Type             | Null | Key | Default | Extra          |
// +----------------------+------------------+------+-----+---------+----------------+
// | id                   | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | uic                  | varchar(255)     | NO   |     |         |                |
// | url                  | varchar(255)     | NO   |     |         |                |
// | callback             | tinyint(4)       | NO   |     | 0       |                |
// | before_callback_sms  | tinyint(4)       | NO   |     | 0       |                |
// | before_callback_mail | tinyint(4)       | NO   |     | 0       |                |
// | after_callback_sms   | tinyint(4)       | NO   |     | 0       |                |
// | after_callback_mail  | tinyint(4)       | NO   |     | 0       |                |
// +----------------------+------------------+------+-----+---------+----------------+

// Action 模型定义和映射
type Action struct {
	ID                 int64  `json:"id"                   gorm:"column:id;type:int;auto_increment;not null;primary_key"`
	UIC                string `json:"uic"                  gorm:"column:uic"`
	URL                string `json:"url"                  gorm:"column:url"`
	Callback           int    `json:"callback"             gorm:"column:callback"`
	BeforeCallbackSMS  int    `json:"before_callback_sms"  gorm:"column:before_callback_sms"`
	BeforeCallbackMail int    `json:"before_callback_mail" gorm:"column:before_callback_mail"`
	AfterCallbackSMS   int    `json:"after_callback_sms"   gorm:"column:after_callback_sms"`
	AfterCallbackMail  int    `json:"after_callback_mail"  gorm:"column:after_callback_mail"`
}

// TableName 结构体映射到的物理表名称
func (m Action) TableName() string {
	return "actions"
}
