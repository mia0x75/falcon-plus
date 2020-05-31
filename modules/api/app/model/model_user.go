package model

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/api/g"
)

// +-----------+------------------+------+-----+------------------+----------------+
// | Field     | Type             | Null | Key | Default          | Extra          |
// +-----------+------------------+------+-----+------------------+----------------+
// | id        | int(10) unsigned | NO   | PRI | NULL             | auto_increment |
// | name      | varchar(50)      | NO   | UNI | NULL             |                |
// | passwd    | char(64)         | NO   |     |                  |                |
// | cnname    | varchar(15)      | NO   |     |                  |                |
// | email     | varchar(75)      | NO   |     |                  |                |
// | phone     | varchar(16)      | NO   |     |                  |                |
// | im        | varchar(50)      | NO   |     |                  |                |
// | role      | tinyint(4)       | NO   |     | 0                |                |
// | creator   | int(10) unsigned | NO   |     | 0                |                |
// | create_at | int(10) unsigned | NO   |     | unix_timestamp() |                |
// +-----------+------------------+------+-----+------------------+----------------+

// User 模型定义和映射
type User struct {
	ID       int64  `json:"id"        gorm:"column:id;type:int;auto_increment;primary_key;not null"`
	Name     string `json:"name"      gorm:"column:name;unique:unique_1;size:50;type:varchar(50);not null"`
	Passwd   string `json:"-"         gorm:"column:passwd;not null;type:char(64);size:64"`
	Cnname   string `json:"cnname"    gorm:"column:cnname;size:15;not null;type:varchar(15)"`
	Email    string `json:"email"     gorm:"column:email;size:75;type:varchar(75);not null"`
	Phone    string `json:"phone"     gorm:"column:phone;size:16;type:varchar(16);not null"`
	IM       string `json:"im"        gorm:"column:im;size:50;type:varchar(50);not null"`
	Role     int    `json:"role"      gorm:"column:role;type:tinyint;not null"`
	Creator  int64  `json:"creator"   gorm:"column:creator;type:int;not null"`
	CreateAt int64  `json:"create_at" gorm:"column:create_at;type:timestamp;not null"`
}

// TableName 结构体映射到的物理表名称
func (m User) TableName() string {
	return "users"
}

// BeforeCreate 插入数据前保证数据的完整性
func (m User) BeforeCreate() (err error) {
	m.CreateAt = time.Now().Unix()
	return
}

func skipAccessControll() bool {
	return !g.Config().AccessControl
}

// IsAdmin 用户是否是管理员
func (m User) IsAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if m.Role == 2 || m.Role == 1 { // TODO: 硬编码
		return true
	}
	return false
}

// IsSuperAdmin 用户是否是超级管理员
func (m User) IsSuperAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if m.Role == 2 { // TODO: 硬编码
		return true
	}
	return false
}
