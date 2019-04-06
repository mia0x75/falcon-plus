package uic

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var db g.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = g.Con()

	u := r.Group("/api/v1/user")
	u.GET("/auth_session", AuthSession)
	u.POST("/login", Login)
	u.GET("/logout", Logout)
	u.POST("/create", CreateUser) // /create -> /

	a := r.Group("/api/v1/user")
	a.Use(utils.AuthSessionMidd)
	a.GET("/current", UserInfo) // /current -> /
	a.GET("/u/:uid", GetUser)
	a.GET("/name/:user_name", GetUserByName)
	a.PUT("/update", UpdateCurrentUser) // /update -> /
	a.PUT("/cgpasswd", ChangePassword)
	a.GET("/users", UserList)
	a.GET("/u/:uid/in_teams", IsUserInTeams)
	a.GET("/u/:uid/teams", GetUserTeams)

	m := r.Group("/api/v1/admin")
	m.Use(utils.AuthSessionMidd)
	m.PUT("/change_user_role", ChangeRoleOfUser)
	m.PUT("/change_user_passwd", AdminChangePassword)
	m.PUT("/change_user_profile", AdminChangeUserProfile)
	m.DELETE("/delete_user", AdminUserDelete)
	m.POST("/login", AdminLogin) // utils.AuthSessionMidd Bug?

	//team
	t := r.Group("/api/v1/team")
	t.Use(utils.AuthSessionMidd)
	t.GET("/", Teams)
	t.GET("/t/:team_id", GetTeam)
	t.GET("/name/:team_name", GetTeamByName)
	t.POST("/", CreateTeam)
	t.PUT("/", UpdateTeam)
	t.POST("/team/user", AddTeamUsers)
	t.DELETE("/:team_id", DeleteTeam)
}
