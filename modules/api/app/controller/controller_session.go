package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
)

// APILoginInput TODO:
type APILoginInput struct {
	Name     string `json:"name"     form:"name"     binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// APIAdminLoginInput TODO:
type APIAdminLoginInput struct {
	Name string `json:"name" form:"name" binding:"required"`
}

// Login TODO:
func Login(c *gin.Context) {
	inputs := APILoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	name := inputs.Name
	password := inputs.Password

	user := model.User{
		Name: name,
	}
	if err := db.Where(&user).Find(&user).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	switch {
	case user.ID == 0:
		h.JSONR(c, h.HTTPBadRequest, "no such user")
		return
	case bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(password)) != nil:
		h.JSONR(c, h.HTTPBadRequest, "password error")
		return
	}
	session := model.Session{}
	if err := db.Where("user_id = ?", user.ID).First(&session).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			h.InternelError(c, "retrieving data", err)
			return
		}
	}
	if session.ID == 0 {
		session.Sign = utils.GenerateUUID()
		session.Expire = int(time.Now().Unix()) + 3600*24*30
		session.UserID = user.ID
		db.Create(&session)
	}
	log.Debugf("[D] Session: %v", session)
	resp := map[string]interface{}{
		"sign":  session.Sign,
		"name":  user.Name,
		"admin": user.IsAdmin(),
	}
	h.JSONR(c, resp)
}

// AdminLogin TODO:
func AdminLogin(c *gin.Context) {
	inputs := APIAdminLoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.Name == inputs.Name {
			return true
		}
		return false
	})
	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (name = %s) does not exist", inputs.Name))
		return
	}
	if user.Role >= ctx.Role {
		// TODO:
		h.JSONR(c, h.HTTPBadRequest, "API_USER not admin, no permissions can do this")
		return
	}
	session := model.Session{}
	if err := db.Where("user_id = ?", user.ID).Find(&session).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			h.InternelError(c, "retrieving data", err)
			return
		}
	}
	if session.ID == 0 {
		session.Sign = utils.GenerateUUID()
		session.Expire = int(time.Now().Unix()) + 3600*24*30
		session.UserID = user.ID
		if err := db.Create(&session).Error; err != nil {
			h.InternelError(c, "creating data", err)
			return
		}
	}
	log.Debugf("[D] Session: %v", session)
	resp := map[string]interface{}{
		"sign":  session.Sign,
		"name":  user.Name,
		"admin": user.IsAdmin(),
	}
	h.JSONR(c, resp)
}

// Logout TODO:
func Logout(c *gin.Context) {
	wsession, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	session := model.Session{}
	user := model.User{}
	if err := db.Where(model.User{Name: wsession.Name}).Find(&user).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	if err := db.Where("sign = ? AND user_id = ?", wsession.Sign, user.ID).Find(&session).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}

	if session.ID == 0 {
		h.JSONR(c, h.HTTPBadRequest, "not found this kind of session in database.")
		return
	}
	if err := db.Delete(&session).Error; err != nil {
		h.InternelError(c, "deleting data", err)
	}
	h.JSONR(c, "logout successful")
}

// AuthSession TODO:
func AuthSession(c *gin.Context) {
	auth, err := h.SessionChecking(c)
	if err != nil || auth != true {
		h.JSONR(c, h.HTTPUnauthorized, err)
		return
	}
	h.JSONR(c, "session is valid!")
}

// CreateRoot TODO:
func CreateRoot(c *gin.Context) {
	password := c.DefaultQuery("password", "")
	if password == "" {
		h.JSONR(c, h.HTTPBadRequest, "password is empty, please check it")
		return
	}
	password = utils.HashIt(password)
	user := model.User{
		Name:   "root",
		Passwd: password,
	}
	if err := db.Create(&user).Error; err != nil {
		h.InternelError(c, "creating data", err)
		return
	}
	h.JSONR(c, "root created!")
}
