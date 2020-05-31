package controller

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/cache"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

// APIUserInput 创建用户的结构体
type APIUserInput struct {
	Name   string `json:"name"     binding:"required"`
	Cnname string `json:"cnname"   binding:"required"`
	Passwd string `json:"password" binding:"required"`
	Email  string `json:"email"    binding:"required"`
	Phone  string `json:"phone"`
	IM     string `json:"im"`
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var inputs APIUserInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if utils.HasDangerousCharacters(inputs.Cnname) {
		h.JSONR(c, h.HTTPBadRequest, "name pattern is invalid")
		return
	}

	// When sign is disabled, only admin user can create user
	if g.Config().SignupDisable {
		ctx, err := h.GetUser(c)
		if err != nil {
			h.JSONR(c, h.HTTPExpectationFailed, err)
			return
		}
		if !ctx.IsAdmin() {
			h.JSONR(c, h.HTTPBadRequest, "Only admin can create new users")
			return
		}
	}
	user := model.User{
		Name:   inputs.Name,
		Cnname: inputs.Cnname,
		Email:  inputs.Email,
		Phone:  inputs.Phone,
		IM:     inputs.IM,
	}
	if cache.UsersMap.Has(user) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (name = %s) already exists", inputs.Name))
		return
	}
	user.Passwd = utils.HashIt(inputs.Passwd)
	// for create a root user during the first time
	if inputs.Name == "root" {
		user.Role = 2
	}

	if err := db.Create(&user).Error; err != nil {
		h.InternelError(c, "creating data", err)
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
	go cache.UsersMap.Init()
}

// APIUserUpdateInput TODO:
type APIUserUpdateInput struct {
	Cnname string `json:"cnname" binding:"required"`
	Email  string `json:"email"  binding:"required"`
	Phone  string `json:"phone"`
	IM     string `json:"im"`
}

// UpdateCurrentUser update current user profile
func UpdateCurrentUser(c *gin.Context) {
	var inputs APIUserUpdateInput
	err := c.Bind(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	case utils.HasDangerousCharacters(inputs.Cnname):
		h.JSONR(c, h.HTTPBadRequest, "name pattern is invalid")
		return
	}
	session, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}
	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.Name == session.Name {
			return true
		}
		return false
	})
	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, "name is not existing")
		return
	}
	update := map[string]interface{}{
		"Cnname": inputs.Cnname,
		"Email":  inputs.Email,
		"Phone":  inputs.Phone,
		"IM":     inputs.IM,
	}
	if err := db.Model(&user).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("User (id = %d) updated", user.ID))
}

// APICgPassedInput TODO:
type APICgPassedInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword TODO:
func ChangePassword(c *gin.Context) {
	var inputs APICgPassedInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
	}
	session, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.Name == session.Name {
			return true
		}
		return false
	})
	if bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(inputs.OldPassword)) != nil {
		h.JSONR(c, h.HTTPBadRequest, "oldPassword is not match current one")
		return
	}

	user.Passwd = utils.HashIt(inputs.NewPassword)
	if err := db.Save(&user).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, "Password updated")
}

// UserInfo 获取当前用户信息
func UserInfo(c *gin.Context) {
	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}
	h.JSONR(c, ctx)
}

// GetUser 根据传入的用户ID获取用户信息
func GetUser(c *gin.Context) {
	inputUserID := c.Params.ByName("id")
	if inputUserID == "" {
		log.Debug("[D] parameter `id` for user is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for user is missing")
		return
	}
	userID, err := strconv.Atoi(inputUserID)
	if err != nil || userID <= 0 {
		log.Debugf("[D] parameter `id` for user is invalid, value = %v", inputUserID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for user is invalid, value = %v", inputUserID))
		return
	}
	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == int64(userID) {
			return true
		}
		return false
	})
	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (id = %d) does not exist", userID))
		return
	}
	h.JSONR(c, user)
}

// UpdateUser TODO:
func UpdateUser(c *gin.Context) {
	inputUserID := c.Params.ByName("id")
	if inputUserID == "" {
		log.Debug("[D] parameter `id` for user is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for user is missing")
		return
	}
	userID, err := strconv.Atoi(inputUserID)
	if err != nil || userID <= 0 {
		log.Debugf("[D] parameter `id` for user is invalid, value = %v", inputUserID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for user is invalid, value = %v", inputUserID))
		return
	}

	var inputs APIUserUpdateInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}
	if utils.HasDangerousCharacters(inputs.Cnname) {
		h.JSONR(c, h.HTTPBadRequest, "name pattern is invalid")
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == int64(userID) {
			return true
		}
		return false
	})
	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (id = %d) does not exist", userID))
		return
	}

	update := map[string]interface{}{
		"Cnname": inputs.Cnname,
		"Email":  inputs.Email,
		"Phone":  inputs.Phone,
		"IM":     inputs.IM,
	}

	if err := db.Model(&user).Update(update).Error; err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("User (id = %d) updated", userID))
	go cache.UsersMap.Init()
}

// GetUserByName 通过登录名获取用户信息
func GetUserByName(c *gin.Context) {
	name := c.Params.ByName("name")
	if name == "" {
		log.Debug("[D] parameter `name` for user is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `name` for user is missing")
		return
	}
	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.Name == name {
			return true
		}
		return false
	})
	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (name = %s) does not exist", name))
		return
	}
	h.JSONR(c, user)
}

// IsUserInTeams TODO:
func IsUserInTeams(c *gin.Context) {
	inputUserID := c.Params.ByName("id")
	if inputUserID == "" {
		log.Debug("[D] parameter `id` for user is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `id` for user is missing")
		return
	}
	userID, err := strconv.Atoi(inputUserID)
	if err != nil || userID <= 0 {
		log.Debugf("[D] parameter `id` for user is invalid, value = %v", inputUserID)
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("Parameter `id` for user is invalid, value = %v", inputUserID))
		return
	}

	inputNames := c.DefaultQuery("names", "")
	if inputNames == "" {
		log.Debug("[D] parameter `names` for team is missing")
		h.JSONR(c, h.HTTPBadRequest, "Parameter `names` for team is missing")
		return
	}
	names := strings.Split(inputNames, ",")

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == int64(userID) {
			return true
		}
		return false
	})

	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("user (id = %d) does not exist", userID))
		return
	}

	teams := cache.TeamsMap.Filter(func(elem *model.Team) bool {
		for _, name := range names {
			if elem.Name == name {
				return true
			}
		}
		return false
	})

	if teams == nil {
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("No matched team(s) found for: %s", inputNames))
		return
	}

	edge := cache.EdgesMap.Any(func(elem *model.Edge) bool {
		for _, t := range teams {
			if elem.AncestorID == t.ID && elem.DescendantID == int64(userID) && elem.Type == 1 {
				return true
			}
		}
		return false
	})

	resp := "true"
	if edge == nil {
		resp = "false"
	}

	h.JSONR(c, resp)
}

// GetUserTeams 获取用户所属的团队信息
func GetUserTeams(c *gin.Context) {
	inputUserID := c.Params.ByName("id")
	if inputUserID == "" {
		h.JSONR(c, h.HTTPBadRequest, "user id is missing")
		return
	}
	userID, err := strconv.Atoi(inputUserID)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err)
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == int64(userID) {
			return true
		}
		return false
	})

	if user == nil {
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("user (id = %d) does not exist", userID))
		return
	}

	edges := cache.EdgesMap.Filter(func(elem *model.Edge) bool {
		if elem.DescendantID == int64(userID) && elem.Type == 1 {
			return true
		}
		return false
	})

	teams := cache.TeamsMap.Filter(func(elem *model.Team) bool {
		for _, r := range edges {
			if elem.ID == r.AncestorID {
				return true
			}
		}
		return false
	})
	resp := map[string]interface{}{
		"teams": teams,
	}
	h.JSONR(c, resp)
}

// APIAdminChangeUserProfileInput 管理员更改用户信息时对应的数据结构
type APIAdminChangeUserProfileInput struct {
	UserID int64  `json:"userID" binding:"required"`
	Cnname string `json:"cnname" binding:"required"`
	Email  string `json:"email"  binding:"required"`
	Phone  string `json:"phone"`
	IM     string `json:"im"`
}

// AdminChangeUserProfile 管理员更改用户信息
func AdminChangeUserProfile(c *gin.Context) {
	var inputs APIAdminChangeUserProfileInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, fmt.Sprintf("An error occurred while parsing input parameter(s), error: %v", err))
		return
	}

	ctx, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.HTTPExpectationFailed, err)
		return
	}
	if !ctx.IsAdmin() {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	if !cache.UsersMap.Has(model.User{ID: inputs.UserID}) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (id = %d) does not exist", inputs.UserID))
		return
	}

	update := map[string]interface{}{
		"ID":     inputs.UserID,
		"Cnname": inputs.Cnname,
		"Email":  inputs.Email,
		"Phone":  inputs.Phone,
		"IM":     inputs.IM,
	}
	if err := db.Model(model.User{}).Where("id = ?", inputs.UserID).Update(update).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, fmt.Sprintf("User (id = %d) updated", inputs.UserID))
	go cache.UsersMap.Init()
}

// APIAdminUserDeleteInput 管理员删除用户对应的数据结构
type APIAdminUserDeleteInput struct {
	UserID int64 `json:"userID" binding:"required"`
}

// AdminUserDelete 管理员删除用户
func AdminUserDelete(c *gin.Context) {
	var inputs APIAdminUserDeleteInput
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
	if !ctx.IsAdmin() {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	if !cache.UsersMap.Has(model.User{ID: inputs.UserID}) {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (id = %d) does not exist", inputs.UserID))
		return
	}
	// TODO: Start transaction to perform:
	// 1. Team -> User ??
	// 2. Table.Creator == inputs.UserID ??
	if err := db.Where("id = ? and role <= ?", inputs.UserID, ctx.Role).Delete(&model.User{}).Error; err != nil {
		h.InternelError(c, "deleting data", err)
		return
	} else if db.RowsAffected == 0 {
		h.JSONR(c, h.HTTPExpectationFailed, "you have no such permission or sth goes wrong")
		return
	}
	h.JSONR(c, fmt.Sprintf("User (id = %d) deleted", inputs.UserID))
	go cache.UsersMap.Init()
}

// APIAdminChangePassword 管理员更改用户密码对应的数据结构
type APIAdminChangePassword struct {
	UserID int64  `json:"userID"   binding:"required"`
	Passwd string `json:"password" binding:"required"`
}

// AdminChangePassword 管理员更改用户密码
func AdminChangePassword(c *gin.Context) {
	var inputs APIAdminChangePassword
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
	if !ctx.IsAdmin() {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == inputs.UserID {
			return true
		}
		return false
	})

	if user == nil {
		h.JSONR(c, h.HTTPBadRequest, fmt.Sprintf("User (id = %d) does not exist", inputs.UserID))
		return
	}

	user.Passwd = utils.HashIt(inputs.Passwd)
	if err := db.Save(user).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}
	h.JSONR(c, "Password updated")
}

// UserList TODO:
func UserList(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	inputPage := c.DefaultQuery("page", "")
	inputLimit := c.DefaultQuery("limit", "")
	page, limit, err = h.PageParser(inputPage, inputLimit)
	if err != nil {
		h.JSONR(c, h.HTTPBadRequest, err.Error())
		return
	}
	q := c.DefaultQuery("q", ".+")
	q = template.HTMLEscapeString(q)
	var user []model.User
	if limit != -1 && page != -1 {
		err = db.Raw("SELECT * FROM users WHERE name REGEXP ? LIMIT ?,?", q, page, limit).Scan(&user).Error
	} else {
		err = db.Where("name REGEXP ?", q).Find(&user).Error
	}
	if err != nil {
		h.InternelError(c, "retrieving data", err)
		return
	}
	h.JSONR(c, user)
}

// APIRoleUpdate TODO:
type APIRoleUpdate struct {
	UserID int64  `json:"userID" binding:"required"`
	Admin  string `json:"admin"  binding:"required"`
}

// ChangeRoleOfUser TODO:
func ChangeRoleOfUser(c *gin.Context) {
	var inputs APIRoleUpdate
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
	if !ctx.IsAdmin() {
		h.JSONR(c, h.HTTPBadRequest, "Permission denied")
		return
	}

	user := cache.UsersMap.Any(func(elem *model.User) bool {
		if elem.ID == inputs.UserID {
			return true
		}
		return false
	})
	switch inputs.Admin {
	case "yes":
		user.Role = 1
	case "no":
		user.Role = 0
	}
	if err := db.Save(user).Error; err != nil {
		h.InternelError(c, "updating data", err)
		return
	}

	h.JSONR(c, "Role changed")
	go cache.UsersMap.Init()
}
