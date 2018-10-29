package portal

import (
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | Field          | Type             | Null | Key | Default           | Extra                       |
// +----------------+------------------+------+-----+-------------------+-----------------------------+
// | id             | int(11)          | NO   | PRI | NULL              | auto_increment              |
// | hostname       | varchar(255)     | NO   | UNI |                   |                             |
// | ip             | varchar(16)      | NO   |     |                   |                             |
// | agent_version  | varchar(16)      | NO   |     |                   |                             |
// | plugin_version | varchar(128)     | NO   |     |                   |                             |
// | maintain_begin | int(10) unsigned | NO   |     | 0                 |                             |
// | maintain_end   | int(10) unsigned | NO   |     | 0                 |                             |
// | update_at      | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +----------------+------------------+------+-----+-------------------+-----------------------------+

type Host struct {
	ID            int64  `json:"id" gorm:"column:id"`
	Hostname      string `json:"hostname" gorm:"column:hostname"`
	Ip            string `json:"ip" gorm:"column:ip"`
	AgentVersion  string `json:"agent_version"  gorm:"column:agent_version"`
	PluginVersion string `json:"plugin_version"  gorm:"column:plugin_version"`
	MaintainBegin int64  `json:"maintain_begin"  gorm:"column:maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"  gorm:"column:maintain_end"`
}

func (this Host) TableName() string {
	return "host"
}

func (this Host) Existing() (int64, bool) {
	db := g.Con()
	db.Falcon.Table(this.TableName()).Where("hostname = ?", this.Hostname).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}

func (this Host) RelatedGrp() (Grps []HostGroup) {
	db := g.Con()
	grpHost := []GrpHost{}
	db.Falcon.Select("grp_id").Where("host_id = ?", this.ID).Find(&grpHost)
	tids := []int64{}
	for _, t := range grpHost {
		tids = append(tids, t.GrpID)
	}
	Grps = []HostGroup{}
	if len(tids) > 0 {
		db.Falcon.Where("id in (?)", tids).Find(&Grps)
	}
	return
}

func (this Host) RelatedTpl() (tpls []Template) {
	db := g.Con()
	grps := this.RelatedGrp()
	gids := []int64{}
	for _, g := range grps {
		gids = append(gids, g.ID)
	}
	grpTpls := []GrpTpl{}
	if len(gids) > 0 {
		db.Falcon.Select("tpl_id").Where("grp_id in (?)", gids).Find(&grpTpls)
	}
	tids := []int64{}
	for _, t := range grpTpls {
		tids = append(tids, t.TplID)
	}
	tpls = []Template{}
	if len(tids) > 0 {
		db.Falcon.Where("id in (?)", tids).Find(&tpls)
	}
	return
}
