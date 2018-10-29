package uic

import (
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

type Team struct {
	ID      int64  `json:"id,"`
	Name    string `json:"name"`
	Resume  string `json:"resume"`
	Creator int64  `json:"creator"`
}

func (this Team) TableName() string {
	return "team"
}

func (this Team) Members() (users []User, err error) {
	db := g.Con()
	var tmapping []RelTeamUser
	if dt := db.Uic.Where("tid = ?", this.ID).Find(&tmapping); dt.Error != nil {
		err = dt.Error
		return
	}
	users = []User{}
	var uids []int64
	for _, t := range tmapping {
		uids = append(uids, t.Uid)
	}

	if len(uids) > 0 {
		if dt := db.Uic.Select("name, id, cnname").Where("id in (?)", uids).Find(&users); dt.Error != nil {
			err = dt.Error
			return
		}
	}
	return
}

func (this Team) GetCreatorName() (userName string, err error) {
	userName = "unknown"
	db := g.Con()
	user := User{ID: this.Creator}
	if dt := db.Uic.Find(&user); dt.Error != nil {
		err = dt.Error
	} else {
		userName = user.Name
	}
	return
}
