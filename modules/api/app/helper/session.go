package helper

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

type WebSession struct {
	Name string
	Sig  string
}

func GetSession(c *gin.Context) (session WebSession, err error) {
	var name, sig string
	apiToken := c.Request.Header.Get("Apitoken")
	if apiToken == "" {
		err = errors.New("token key is not set")
		return
	}
	log.Debugf("[D] header: %v, apiToken: %v", c.Request.Header, apiToken)
	var websession WebSession
	err = json.Unmarshal([]byte(apiToken), &websession)
	if err != nil {
		return
	}
	name = websession.Name
	log.Debugf("[D] session got name: %s", name)
	if name == "" {
		err = errors.New("token key:name is empty")
		return
	}
	sig = websession.Sig
	log.Debugf("[D] session got sig: %s", sig)
	if sig == "" {
		err = errors.New("token key:sig is empty")
		return
	}
	if err != nil {
		return
	}
	session = WebSession{name, sig}
	return
}

func SessionChecking(c *gin.Context) (auth bool, err error) {
	auth = false
	var websessio WebSession
	websessio, err = GetSession(c)
	if err != nil {
		return
	}

	//default_token used in server side access
	default_token := g.Config().DefaultToken
	if default_token != "" && websessio.Sig == default_token {
		auth = true
		return
	}

	db := g.Con().Uic
	var user uic.User
	db.Where("name = ?", websessio.Name).Find(&user)
	if user.ID == 0 {
		err = errors.New("not found this user")
		return
	}
	var session uic.Session
	db.Table("session").Where("sig = ? and uid = ?", websessio.Sig, user.ID).Scan(&session)
	if session.ID == 0 {
		err = errors.New("session not found")
		return
	} else {
		auth = true
	}
	return
}

func GetUser(c *gin.Context) (user uic.User, err error) {
	db := g.Con().Uic
	websession, getserr := GetSession(c)
	if getserr != nil {
		err = getserr
		return
	}
	user = uic.User{
		Name: websession.Name,
	}
	dt := db.Table("user").Where(&user).Find(&user)
	err = dt.Error
	return
}
