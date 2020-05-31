package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/set"
	"github.com/toolkits/net/httplib"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// APIGetTeamOutput TODO:
type APIGetTeamOutput struct {
	model.Team
	Users       []*model.User `json:"users"`
	TeamCreator string        `json:"teamCreator"`
}

// UsersCache TODO:
type UsersCache struct {
	sync.RWMutex
	M map[string][]*model.User
}

// Users TODO:
var Users = &UsersCache{M: make(map[string][]*model.User)}

// Get TODO:
func (s *UsersCache) Get(team string) []*model.User {
	s.RLock()
	defer s.RUnlock()
	val, exists := s.M[team]
	if !exists {
		return nil
	}

	return val
}

// Set TODO:
func (s *UsersCache) Set(team string, users []*model.User) {
	s.Lock()
	defer s.Unlock()
	s.M[team] = users
}

// UsersOf TODO:
func UsersOf(team string) []*model.User {
	users := CurlUic(team)

	if users != nil {
		Users.Set(team, users)
	} else {
		users = Users.Get(team)
	}

	return users
}

// GetUsers TODO:
func GetUsers(teams string) map[string]*model.User {
	userMap := make(map[string]*model.User)
	arr := strings.Split(teams, ",")
	for _, team := range arr {
		if team == "" {
			continue
		}

		users := UsersOf(team)
		if users == nil {
			continue
		}

		for _, user := range users {
			userMap[user.Name] = user
		}
	}
	return userMap
}

// ParseTeams return phones, emails, IM
func ParseTeams(teams string) ([]string, []string, []string) {
	if teams == "" {
		return []string{}, []string{}, []string{}
	}

	userMap := GetUsers(teams)
	phoneSet := set.NewStringSet()
	mailSet := set.NewStringSet()
	imSet := set.NewStringSet()
	for _, user := range userMap {
		if user.Phone != "" {
			phoneSet.Add(user.Phone)
		}
		if user.Email != "" {
			mailSet.Add(user.Email)
		}
		if user.IM != "" {
			imSet.Add(user.IM)
		}
	}
	return phoneSet.ToSlice(), mailSet.ToSlice(), imSet.ToSlice()
}

// CurlUic TODO:
func CurlUic(team string) []*model.User {
	if team == "" {
		return []*model.User{}
	}

	uri := fmt.Sprintf("%s/api/v1/team/name/%s", g.Config().API.API, team)
	req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  g.Config().API.Token,
	})
	req.Header("X-Falcon-Token", string(token))

	var teamUsers APIGetTeamOutput
	err := req.ToJson(&teamUsers)
	if err != nil {
		log.Errorf("[E] curl %s fail: %v", uri, err)
		return nil
	}

	return teamUsers.Users
}
