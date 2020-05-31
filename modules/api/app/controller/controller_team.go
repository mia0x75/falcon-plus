package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// APITeamsOutput TODO:
type APITeamsOutput struct {
	Team        model.Team    `json:"team"`
	TeamCreator string        `json:"creator"`
	Users       []*model.User `json:"users"`
}

// Teams support root as admin
func Teams(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	query := c.DefaultQuery("q", ".+")
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}
	teams := []model.Team{}
	if ctx.IsAdmin() {
		if limit != -1 && page != -1 {
			err = db.Raw("SELECT * FROM teams WHERE name REGEXP ? LIMIT ?, ?", query, page, limit).
				Scan(&teams).
				Error
		} else {
			err = db.Where("name REGEXP ?", query).Find(&teams).Error
		}
	} else {
		// team creator and team member can manage the team
		err = db.Raw(`SELECT t.* FROM teams as t, edges as l
			WHERE t.name REGEXP ? AND t.id = l.ancestor_id AND l.descendant_id = ? AND l.type = 1
			UNION SELECT * FROM teams WHERE name REGEXP ? AND creator = ?`,
			query, ctx.ID, query, ctx.ID).
			Scan(&teams).
			Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	f := func(t model.Team) []*model.User {
		edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
			if elem.AncestorID == t.ID && elem.Type == 1 {
				return true
			}
			return false
		})
		if edges == nil {
			return nil
		}
		users := cache.UsersMap.Filter(func(elem *model.User) bool {
			for _, l := range edges {
				if elem.ID == l.DescendantID {
					return true
				}
			}
			return false
		})
		return users
	}
	resp := []APITeamsOutput{}
	for _, t := range teams {
		output := APITeamsOutput{
			Team: t,
		}
		output.Users = f(t)
		creator := cache.UsersMap.Any(func(elem *model.User) bool {
			if elem.ID == t.Creator {
				return true
			}
			return false
		})
		if creator != nil {
			output.TeamCreator = creator.Name
		} else {
			output.TeamCreator = "unknown"
		}

		resp = append(resp, output)
	}
	h.JSONR(c, resp)
}

// APICreateTeamInput TODO:
type APICreateTeamInput struct {
	Name    string  `json:"name"   binding:"required"`
	Resume  string  `json:"resume"`
	UserIDs []int64 `json:"users"`
}

// CreateTeam every user can create a team
func CreateTeam(c *gin.Context) {
	var inputs APICreateTeamInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	team := model.Team{
		Name:    inputs.Name,
		Resume:  inputs.Resume,
		Creator: ctx.ID,
	}
	if cache.TeamsMap.Has(team) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Team (name = %s) already exists", inputs.Name))
		return
	}
	tx := db.Begin()
	if err := tx.Create(&team).Error; err != nil {
		h.InternelError(c, "creating data", err)
		tx.Rollback()
		return
	}
	if len(inputs.UserIDs) > 0 {
		for _, userID := range inputs.UserIDs {
			if !cache.UsersMap.Has(model.User{ID: userID}) {
				continue
			}
			edge := &model.Edge{
				AncestorID:   team.ID,
				DescendantID: userID,
				Type:         1,
				Creator:      ctx.ID,
			}
			if err := tx.Create(edge).Error; err != nil {
				h.InternelError(c, "creating data", err)
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Team (id = %d, name = %s) created!", team.ID, team.Name))
	go cache.TeamsMap.Init()
	go cache.EdgesMap.Init()
}

// APIUpdateTeamInput TODO:
type APIUpdateTeamInput struct {
	ID      int    `json:"id"     binding:"required"`
	Resume  string `json:"resume"`
	Name    string `json:"name"`
	UserIDs []int  `json:"users"`
}

// UpdateTeam admin, team creator, team member can manage the team
func UpdateTeam(c *gin.Context) {
	var (
		inputs APIUpdateTeamInput
		team   *model.Team
		err    error
	)
	err = c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	team = cache.TeamsMap.Any(func(elem *model.Team) bool {
		if elem.ID == int64(inputs.ID) {
			return true
		}
		return false
	})
	if team == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Team (id = %d) does not exits", inputs.ID))
		return
	}
	// 如果不是管理员也不是创建者
	if !ctx.IsAdmin() && team.Creator != ctx.ID {
		edge := cache.EdgesMap.Any(func(elem *model.Edge) bool {
			if elem.AncestorID == team.ID && elem.DescendantID == ctx.ID && elem.Type == 1 {
				return true
			}
			return true
		})
		if edge == nil {
			// 同时也不是团队成员
			team = nil
		}
	}
	if team == nil {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	update := map[string]interface{}{
		"Name":   inputs.Name,
		"Resume": inputs.Resume,
	}
	if err := db.Model(team).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}

	// TODO:
	if err := bindUsers(db, inputs.ID, inputs.UserIDs); err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("Team (id = %d) updated!", inputs.ID))
	go cache.TeamsMap.Init()
}

// APIAddTeamUsers TODO:
type APIAddTeamUsers struct {
	TeamID int      `json:"teamID" binding:"required"`
	Users  []string `json:"users"  binding:"required"`
}

// AddTeamUsers admin, team creator, team member can mangage the team
func AddTeamUsers(c *gin.Context) {
	var (
		inputs APIAddTeamUsers
		team   *model.Team
		err    error
	)
	err = c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	team = cache.TeamsMap.Any(func(elem *model.Team) bool {
		if elem.ID == int64(inputs.TeamID) {
			return true
		}
		return false
	})
	if team == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Team (id = %d) does not exits", inputs.TeamID))
		return
	}
	// 如果不是管理员也不是创建者
	if !ctx.IsAdmin() && team.Creator != ctx.ID {
		edge := cache.EdgesMap.Any(func(elem *model.Edge) bool {
			if elem.AncestorID == team.ID && elem.DescendantID == ctx.ID && elem.Type == 1 {
				return true
			}
			return true
		})
		if edge == nil {
			// 同时也不是团队成员
			team = nil
		}
	}
	if team == nil {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}
	users := cache.UsersMap.Filter(func(elem *model.User) bool {
		for _, n := range inputs.Users {
			if elem.Name == n {
				return true
			}
		}
		return false
	})
	if users == nil {
		h.JSONR(c, h.HTTPBadRequest, "empty users")
		return
	}
	for _, u := range users {
		edge := model.Edge{
			AncestorID:   int64(inputs.TeamID),
			DescendantID: int64(u.ID),
			Type:         1,
			Creator:      ctx.ID,
		}
		if !cache.EdgesMap.Has(edge) {
			if err := db.Create(&edge).Error; err != nil {
				h.InternelError(c, "creating data", err)
				return
			}
		}
	}
	h.JSONR(c, "add successful")
	go cache.EdgesMap.Init()
}

func bindUsers(db *gorm.DB, tid int, users []int) (err error) {
	if len(users) == 0 {
		return
	}

	// TODO: FIX delete unbind users
	var needDeleteMan []model.Edge
	if err = db.Where("ancestor_id = ? AND NOT (descendant_id IN (?)) AND type = 1", tid, users).Find(&needDeleteMan).Error; err != nil {
		return
	}
	if len(needDeleteMan) != 0 {
		for _, man := range needDeleteMan {
			if err = db.Delete(&man).Error; err != nil {
				return
			}
		}
	}

	// insert bind users
	for _, i := range users {
		edge := model.Edge{
			AncestorID:   int64(tid),
			DescendantID: int64(i),
			Type:         1,
		}
		if err = db.Where(&edge).Find(&edge).Error; err != nil {
			return
		}
		if edge.ID == 0 {
			if err = db.Create(&edge).Error; err != nil {
				return
			}
		} else {
			// if record exist, do next
			continue
		}
	}
	return
}

// DeleteTeam only admin or team creator can delete a team
func DeleteTeam(c *gin.Context) {
	var err error
	inputTeamID := c.Params.ByName("id")
	if inputTeamID == "" {
		log.Debug("[D] parameter `id` for team is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for team is missing")
		return
	}
	teamID, err := strconv.Atoi(inputTeamID)
	if err != nil {
		log.Debugf("[D] parameter `id` for team is invalid, value = %v", inputTeamID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for team is invalid, value = %v", inputTeamID))
		return
	}
	if teamID == 0 {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Invalid team (id = %d)", teamID))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	team := cache.TeamsMap.Any(func(elem *model.Team) bool {
		if elem.ID == int64(teamID) {
			return true
		}
		return false
	})
	if team == nil {
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("Team (id = %d) does not exist", teamID))
		return
	}

	if !ctx.IsAdmin() && team.Creator != ctx.ID {
		h.JSONR(c, h.HTTPExpectationFailed, "Permission denied")
		return
	}
	tx := db.Begin()
	if err := tx.Delete(team).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}

	if err := tx.Where("ancestor_id = ? AND type = 1", teamID).Delete(model.Edge{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("Team (id = %d) deleted", teamID))
	go cache.TeamsMap.Init()
	go cache.EdgesMap.Init()
}

// APIGetTeamOutput TODO:
type APIGetTeamOutput struct {
	*model.Team
	Users []*model.User `json:"users"`
}

// GetTeam 根据主键获取一个团队的信息
func GetTeam(c *gin.Context) {
	inputTeamID := c.Params.ByName("id")
	if inputTeamID == "" {
		log.Debug("[D] parameter `id` for team is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for team is missing")
		return
	}
	teamID, err := strconv.Atoi(inputTeamID)
	if err != nil || teamID <= 0 {
		log.Debugf("[D] parameter `id` for team is invalid, value = %v", inputTeamID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for team is invalid, value = %v", inputTeamID))
		return
	}
	team := cache.TeamsMap.Any(func(elem *model.Team) bool {
		if elem.ID == int64(teamID) {
			return true
		}
		return false
	})
	if team == nil {
		log.Debugf("[D] team (id = %d) does not exist", teamID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Team (id = %d) does not exist", teamID))
		return
	}
	resp := GetTeamInfo(team)
	h.JSONR(c, resp)
}

// GetTeamByName 通过团队名称获取团队信息
func GetTeamByName(c *gin.Context) {
	name := c.Params.ByName("name")
	if name == "" {
		log.Debug("[D] parameter `name` for team is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `name` for team is missing")
		return
	}
	team := cache.TeamsMap.Any(func(elem *model.Team) bool {
		if elem.Name == name {
			return true
		}
		return false
	})
	if team == nil {
		log.Debugf("[D] team (name = %s) does not exist", name)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Team (name = %s) does not exist", name))
		return
	}
	resp := GetTeamInfo(team)
	h.JSONR(c, resp)
}

// GetTeamInfo 获取团队信息
func GetTeamInfo(team *model.Team) (resp APIGetTeamOutput) {
	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.AncestorID == team.ID && elem.Type == 1 {
			return true
		}
		return false
	})
	resp.Team = team
	if edges != nil {
		users := cache.UsersMap.Filter(func(elem *model.User) bool {
			for _, r := range edges {
				if elem.ID == r.DescendantID {
					return true
				}
			}
			return false
		})
		resp.Users = users
	}
	return
}
