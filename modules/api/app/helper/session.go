package helper

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

type WebSession struct {
	Name string
	Sign string
}

func GetSession(c *gin.Context) (session WebSession, err error) {
	var name, sign string
	apiToken := c.Request.Header.Get("X-Falcon-Token")
	if apiToken == "" {
		err = errors.New("token key is not set")
		return
	}
	log.Debugf("[D] header: %v, apiToken: %v", c.Request.Header, apiToken)
	var s WebSession
	err = json.Unmarshal([]byte(apiToken), &s)
	if err != nil {
		return
	}
	name = s.Name
	log.Debugf("[D] session got name: %s", name)
	if name == "" {
		err = errors.New("token key:name is empty")
		return
	}
	sign = s.Sign
	log.Debugf("[D] session got sign: %s", sign)
	if sign == "" {
		err = errors.New("token key:sign is empty")
		return
	}
	if err != nil {
		return
	}
	session = WebSession{name, sign}
	return
}

func SessionChecking(c *gin.Context) (auth bool, err error) {
	auth = false
	paths := []string{
		"/",
		"/health",
		"/version",
		"/workdir",
		"/config",
		"/config/reload",
		"/api/v1/user/auth_session",
		"/api/v1/user/login",
		"/api/v1/user/logout",
		"/api/v1/user/create",
		"/api/v1/grafana",
		"/api/v1/grafana/metrics/find",
		"/api/v1/grafana/render",
		"/api/v1/grafana/render",
	}

	for _, p := range paths {
		if c.FullPath() == p {
			auth = true
			return
		}
	}
	var s WebSession
	s, err = GetSession(c)
	if err != nil {
		return
	}

	// default_token used in server side access
	default_token := g.Config().DefaultToken
	if default_token != "" && s.Sign == default_token {
		auth = true
		return
	}

	db := g.Con()
	var user model.User
	db.Where("name = ?", s.Name).Find(&user)
	if user.ID == 0 {
		err = errors.New("not found this user")
		return
	}
	session := model.Session{}
	db.Table(session.TableName()).Where("sign = ? and user_id = ?", s.Sign, user.ID).Find(&session)
	if session.ID == 0 {
		err = errors.New("session not found")
		return
	} else {
		auth = true
	}
	return
}

func GetUser(c *gin.Context) (user model.User, err error) {
	db := g.Con()
	var s WebSession
	if s, err = GetSession(c); err != nil {
		return
	} else {
		user = model.User{
			Name: s.Name,
		}
		dt := db.Table(user.TableName()).Where(&user).Find(&user)
		err = dt.Error
		return
	}
}
