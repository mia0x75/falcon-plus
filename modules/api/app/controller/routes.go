package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var db *gorm.DB

const (
	TMP_GRAPH_FILED_DELIMITER = "|"
)

// Routes 路由表
func Routes(r *gin.Engine) {
	db = g.Con()
	userRoutes(r)
	adminRoutes(r)
	teamRoutes(r)

	groupRoute(r)
	pluginRoute(r)
	aggregatorRoute(r)
	templateRoute(r)
	hostRoute(r)
	maintainRoute(r)
	strategyRoute(r)
	metricRoute(r)
	endpointRoute(r)
	actionRoute(r)
	nodataRoute(r)
	expressionRoute(r)

	alarmRoutes(r)

	dashboardRoutes(r)

	graphRoutes(r)
	grafanaRoutes(r)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/user/logout                    |                                |
// | GET    | /api/v1/user/auth_session              |                                |
// | POST   | /api/v1/user/login                     |                                |
// | POST   | /api/v1/user/create                    |                                |
// | GET    | /api/v1/user/current                   |                                |
// | GET    | /api/v1/user/id/:id                    |                                |
// | PUT    | /api/v1/user/id/:id                    |                                |
// | GET    | /api/v1/user/name/:name                |                                |
// | PUT    | /api/v1/user/update                    |                                |
// | PUT    | /api/v1/user/cgpasswd                  |                                |
// | GET    | /api/v1/user/users                     |                                |
// | GET    | /api/v1/user/id/:id/in_teams           |                                |
// | GET    | /api/v1/user/id/:id/teams              |                                |
// +--------+----------------------------------------+--------------------------------+
func userRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/user")
	g.GET("/auth_session", AuthSession)
	g.POST("/login", Login)
	g.GET("/logout", Logout)
	g.POST("/create", CreateUser)
	g.GET("/current", UserInfo)
	g.GET("/id/:id", GetUser)
	g.PUT("/id/:id", UpdateUser)
	g.GET("/name/:name", GetUserByName)
	g.PUT("/update", UpdateCurrentUser)
	g.PUT("/cgpasswd", ChangePassword)
	g.GET("/users", UserList)
	g.GET("/id/:id/in_teams", IsUserInTeams)
	g.GET("/id/:id/teams", GetUserTeams)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | PUT    | /api/v1/admin/change_user_role         |                                |
// | PUT    | /api/v1/admin/change_user_passwd       |                                |
// | PUT    | /api/v1/admin/change_user_profile      |                                |
// | DELETE | /api/v1/admin/delete_user              |                                |
// | POST   | /api/v1/admin/login                    |                                |
// +--------+----------------------------------------+--------------------------------+
func adminRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/admin")
	g.PUT("/change_user_role", ChangeRoleOfUser)
	g.PUT("/change_user_passwd", AdminChangePassword)
	g.PUT("/change_user_profile", AdminChangeUserProfile)
	g.DELETE("/delete_user", AdminUserDelete)
	g.POST("/login", AdminLogin) // utils.AuthSessionMidd Bug?
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/team                           |                                |
// | POST   | /api/v1/team                           |                                |
// | PUT    | /api/v1/team                           |                                |
// | GET    | /api/v1/team/id/:id                    |                                |
// | DELETE | /api/v1/team/id/:id                    |                                |
// | GET    | /api/v1/team/name/:name                |                                |
// | POST   | /api/v1/team/team/user                 |                                |
// +--------+----------------------------------------+--------------------------------+
func teamRoutes(r *gin.Engine) {
	// team
	g := r.Group("/api/v1/team")
	g.GET("", Teams)
	g.POST("", CreateTeam)
	g.PUT("", UpdateTeam)
	g.GET("/id/:id", GetTeam)
	g.DELETE("/id/:id", DeleteTeam)
	g.GET("/name/:name", GetTeamByName)
	g.POST("/team/user", AddTeamUsers)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/endpoint/name/:name/group      |                                |
// +--------+----------------------------------------+--------------------------------+
func endpointRoute(r *gin.Engine) {
	g := r.Group("/api/v1/endpoint")
	g.GET("/name/:name/group", GetEndpointRelatedGroups)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | POST   | /api/v1/aggregator                     |                                |
// | PUT    | /api/v1/aggregator                     |                                |
// | GET    | /api/v1/aggregator/id/:id              |                                |
// | DELETE | /api/v1/aggregator/id/:id              |                                |
// +--------+----------------------------------------+--------------------------------+
func aggregatorRoute(r *gin.Engine) {
	g := r.Group("/api/v1/aggregator")
	g.POST("", CreateAggregator)
	g.PUT("", UpdateAggregator)
	g.GET("/id/:id", GetAggregator)
	g.DELETE("/id/:id", DeleteAggregator)
	g.GET("/group/id/:id", GetAggregatorListOfGroup)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | POST   | /api/v1/plugin                         |                                |
// | GET    | /api/v1/plugin/group/id/:id            |                                |
// | DELETE | /api/v1/plugin/id/:id                  |                                |
// +--------+----------------------------------------+--------------------------------+
func pluginRoute(r *gin.Engine) {
	g := r.Group("/api/v1/plugin")
	g.POST("", CreatePlugin)
	g.GET("/group/id/:id", GetPluginOfGroup)
	g.DELETE("/id/:id", DeletePlugin)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/strategy                       |                                |
// | POST   | /api/v1/strategy                       |                                |
// | PUT    | /api/v1/strategy                       |                                |
// | GET    | /api/v1/strategy/id/:id                |                                |
// | DELETE | /api/v1/strategy/id/:id                |                                |
// +--------+----------------------------------------+--------------------------------+
func strategyRoute(r *gin.Engine) {
	g := r.Group("/api/v1/strategy")
	g.GET("", GetStrategies)
	g.POST("", CreateStrategy)
	g.PUT("", UpdateStrategy)
	g.GET("/id/:id", GetStrategy)
	g.DELETE("/id/:id", DeleteStrategy)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/metric/default_list            |                                |
// +--------+----------------------------------------+--------------------------------+
func metricRoute(r *gin.Engine) {
	g := r.Group("/api/v1/metric")
	g.GET("/default_list", MetricQuery)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/template                       |                                |
// | POST   | /api/v1/template                       |                                |
// | PUT    | /api/v1/template                       |                                |
// | POST   | /api/v1/template/action                |                                |
// | PUT    | /api/v1/template/action                |                                |
// | GET    | /api/v1/template/id/:id                |                                |
// | DELETE | /api/v1/template/id/:id                |                                |
// | GET    | /api/v1/template/id/:id/group          |                                |
// +--------+----------------------------------------+--------------------------------+
func templateRoute(r *gin.Engine) {
	g := r.Group("/api/v1/template")
	g.GET("", GetTemplates)
	g.POST("", CreateTemplate)
	g.PUT("", UpdateTemplate)

	// TODO:
	// g.GET("/template_simple", GetTemplatesSimple)
	g.POST("/action", CreateActionToTmplate)
	g.PUT("/action", UpdateActionToTmplate)
	g.GET("/id/:id", GetTemplate)
	g.DELETE("/id/:id", DeleteTemplate)
	g.GET("/id/:id/group", GetATemplateHostgroup)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/action/id/:id                  |                                |
// +--------+----------------------------------------+--------------------------------+
func actionRoute(r *gin.Engine) {
	g := r.Group("/api/v1/action")
	g.GET("/id/:id", GetAction)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/host                           |                                |
// | POST   | /api/v1/host/find_by_strategy          |                                |
// | GET    | /api/v1/host/id/:id/template           |                                |
// | GET    | /api/v1/host/id/:id/group              |                                |
// +--------+----------------------------------------+--------------------------------+
func hostRoute(r *gin.Engine) {
	g := r.Group("/api/v1/host")
	g.GET("", GetHosts)
	g.POST("/find_by_strategy", FindByMetric) // TODO: 这个API的URL不太好
	g.GET("/id/:id/template", GetHostRelatedTemplates)
	g.GET("/id/:id/hostgroup", GetHostRelatedGroups)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/maintain                       |                                |
// | POST   | /api/v1/maintain                       |                                |
// | DELETE | /api/v1/maintain                       |                                |
// +--------+----------------------------------------+--------------------------------+
func maintainRoute(r *gin.Engine) {
	g := r.Group("/api/v1/maintain")
	g.GET("", GetMaintain)
	g.POST("", SetMaintain)
	g.DELETE("", UnsetMaintain)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/group                          |                                |
// | POST   | /api/v1/group                          |                                |
// | PUT    | /api/v1/group                          |                                |
// | POST   | /api/v1/group/host                     |                                |
// | PUT    | /api/v1/group/host                     |                                |
// | POST   | /api/v1/group/template                 |                                |
// | PUT    | /api/v1/group/template                 |                                |
// | GET    | /api/v1/group/id/:id                   |                                |
// | DELETE | /api/v1/group/id/:id                   |                                |
// | PATCH  | /api/v1/group/id/:id/host              |                                |
// | GET    | /api/v1/group/id/:id/template          |                                |
// +--------+----------------------------------------+--------------------------------+
func groupRoute(r *gin.Engine) {
	g := r.Group("/api/v1/group")
	g.GET("", GetGroups)
	g.POST("", CreateGroup)
	g.PUT("", PutGroup)
	g.POST("/host", AddHostsToGroup)
	g.PUT("/host", RemoveHostFromGroup)
	g.POST("/template", BindTemplateToGroup)
	g.PUT("/template", UnBindTemplateToGroup)

	g.GET("/id/:id", GetGroup)
	g.DELETE("/id/:id", DeleteGroup)
	g.PATCH("/id/:id/host", PatchHostInGroup)
	g.GET("/id/:id/template", GetTemplateOfGroup)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/nodata                         |                                |
// | POST   | /api/v1/nodata                         |                                |
// | PUT    | /api/v1/nodata                         |                                |
// | GET    | /api/v1/nodata/id/:id                  |                                |
// | DELETE | /api/v1/nodata/id/:id                  |                                |
// +--------+----------------------------------------+--------------------------------+
func nodataRoute(r *gin.Engine) {
	g := r.Group("/api/v1/nodata")
	g.GET("", GetMockcfgList)
	g.POST("", CreateMockcfg)
	g.PUT("", UpdateMockcfg)
	g.GET("/id/:id", GetMockcfg)
	g.DELETE("/id/:id", DeleteMockcfg)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/expression                     |                                |
// | POST   | /api/v1/expression                     |                                |
// | PUT    | /api/v1/expression                     |                                |
// | GET    | /api/v1/expression/id/:id              |                                |
// | DELETE | /api/v1/expression/id/:id              |                                |
// +--------+----------------------------------------+--------------------------------+
func expressionRoute(r *gin.Engine) {
	g := r.Group("/api/v1/expression")
	g.GET("", GetExpressionList)
	g.POST("", CreateExrpession)
	g.PUT("", UpdateExrpession)
	g.GET("/id/:id", GetExpression)
	g.DELETE("/id/:id", DeleteExpression)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/alarm/cases                    |                                |
// | POST   | /api/v1/alarm/cases                    |                                |
// | GET    | /api/v1/alarm/events                   |                                |
// | POST   | /api/v1/alarm/events                   |                                |
// | GET    | /api/v1/alarm/notes                    |                                |
// | POST   | /api/v1/alarm/notes                    |                                |
// +--------+----------------------------------------+--------------------------------+
func alarmRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/alarm")
	g.GET("/cases", Lists)
	g.POST("/cases", Lists)
	g.GET("/events", GetEvents)
	g.POST("/events", GetEvents)
	g.GET("/notes", GetNotesOfAlarm)
	g.POST("/notes", AddNotesToAlarm)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/dashboard/tmpgraph/id/:id      |                                |
// | POST   | /api/v1/dashboard/tmpgraph             |                                |
// | GET    | /api/v1/dashboard/graph/id/:id         |                                |
// | PUT    | /api/v1/dashboard/graph/id/:id         |                                |
// | DELETE | /api/v1/dashboard/graph/id/:id         |                                |
// | POST   | /api/v1/dashboard/graph                |                                |
// | GET    | /api/v1/dashboard/graphs/screen/id/:id |                                |
// | GET    | /api/v1/dashboard/screen/id/:id        |                                |
// | POST   | /api/v1/dashboard/screen               |                                |
// | GET    | /api/v1/dashboard/screens/pid/:pid     |                                |
// | GET    | /api/v1/dashboard/screens              |                                |
// | PUT    | /api/v1/dashboard/screen/id/:id        |                                |
// | DELETE | /api/v1/dashboard/screen/id/:id        |                                |
// +--------+----------------------------------------+--------------------------------+
func dashboardRoutes(r *gin.Engine) {
	db = g.Con()
	g := r.Group("/api/v1/dashboard")
	g.GET("/tmpgraph/id/:id", GetTempGraph)
	g.POST("/tmpgraph", CreateTempGraph)
	g.GET("/graph/id/:id", GetGraph)
	g.PUT("/graph/id/:id", UpdateGraph)
	g.POST("/graph", CreateGraph)
	g.DELETE("/graph/id/:id", DeleteGraph)
	g.GET("/graphs/screen/id/:id", GetGraphsByScreenID)
	g.GET("/screen/id/:id", GetScreen)
	g.POST("/screen", CreateScreen)
	g.GET("/screens/pid/:pid", GetScreensByPid)
	g.GET("/screens", GetScreensAll)
	g.PUT("/screen/id/:id", UpdateScreen)
	g.DELETE("/screen/id/:id", DeleteScreen)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/graph/endpointobj              |                                |
// | GET    | /api/v1/graph/endpoint                 |                                |
// | GET    | /api/v1/graph/endpoint_counter         |                                |
// | POST   | /api/v1/graph/history                  |                                |
// | POST   | /api/v1/graph/lastpoint                |                                |
// | DELETE | /api/v1/graph/endpoint                 |                                |
// | DELETE | /api/v1/graph/counter                  |                                |
// +--------+----------------------------------------+--------------------------------+
func graphRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/graph")
	g.GET("/endpointobj", EndpointObjGet)
	g.GET("/endpoint", EndpointRegexpQuery)
	g.GET("/endpoint_counter", EndpointCounterRegexpQuery)
	g.POST("/history", QueryGraphDrawData)
	g.POST("/lastpoint", QueryGraphLastPoint)
	g.DELETE("/endpoint", DeleteGraphEndpoint)
	g.DELETE("/counter", DeleteGraphCounter)
}

// +--------+----------------------------------------+--------------------------------+
// | Method | Path                                   | Description                    |
// +--------+----------------------------------------+--------------------------------+
// | GET    | /api/v1/grafana                        |                                |
// | GET    | /api/v1/grafana/metrics/find           |                                |
// | POST   | /api/v1/grafana/render                 |                                |
// | GET    | /api/v1/grafana/render                 |                                |
// +--------+----------------------------------------+--------------------------------+
func grafanaRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/grafana")
	g.GET("", GrafanaMainQuery)
	g.GET("/metrics/find", GrafanaMainQuery)
	g.POST("/render", GrafanaRender)
	g.GET("/render", GrafanaRender)
}
