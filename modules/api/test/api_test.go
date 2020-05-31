package test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/g"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	api_v1             = ""
	test_user_name     = "apitest-user1"
	test_user_password = "password"
	test_user_id       = 0
	test_team_name     = "apitest-team1"
	test_team_id       = 0
	root_user_name     = "root"
	root_user_password = "rootpass"
	root_user_id       = 0
	test_group_id      = 0
	test_group_name    = "apitest-group1"
)

func init() {
	g.ParseConfig("../cfg.example.json")
	cu.InitLog(g.Config().Log.Level)
	port := strings.TrimLeft(g.Config().Listen, ":")
	host := "127.0.0.1"
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("[F] %v", err)
	}
	api_v1 = fmt.Sprintf("http://%s:%s/api/v1", host, port)

	if err := g.InitDB(); err != nil {
		log.Fatalf("[F] %v", err)
	}
	init_testing_data()
}

func init_testing_data() {
	password := utils.HashIt(test_user_password)
	user := model.User{
		Name:   test_user_name,
		Passwd: password,
		Cnname: test_user_name,
		Email:  test_user_name + "@test.com",
		Phone:  "1234567890",
		IM:     "hellotest",
	}

	db := g.Con()
	db.Table("users").Where("name = ?", test_user_name).Delete(&model.User{})
	if err := db.Table("users").Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	log.Info("create_user:", test_user_name)
	log.Info("user_id:", user.ID)
	test_user_id = int(user.ID)

	db.Table("users").Where("name = ?", "root").Delete(&model.User{})
	db.Table("teams").Where("name = ?", test_team_name).Delete(&model.Team{})
}

func get_session_token() (string, error) {
	rr := map[string]interface{}{}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"name":     root_user_name,
			"password": root_user_password,
		}).
		SetResult(&rr).
		Post(fmt.Sprintf("%s/user/login", api_v1))

	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", errors.New(resp.String())
	}

	api_token := fmt.Sprintf(`{"name": "%v", "sign": "%v"}`, rr["name"], rr["sign"])
	return api_token, nil
}

func TestUser(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}
	var api_token string = ""

	Convey("Create root user: POST /user/create", t, func() {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"name":     root_user_name,
				"password": root_user_password,
				"email":    "root@test.com",
				"cnname":   "cnroot",
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/create", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
		}
	})

	Convey("Login user: POST /user/login", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"name":     root_user_name,
				"password": root_user_password,
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, root_user_name)
			So((*rr)["sign"], ShouldNotBeBlank)
			So((*rr)["admin"], ShouldBeTrue)
			api_token = fmt.Sprintf(`{"name": "%v", "sign": "%v"}`, (*rr)["name"], (*rr)["sign"])
		}
	})

	Convey("Get user info by name: GET /user/name/:name", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["role"], ShouldEqual, 2)
			So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
		}
	})
	root_user_id := (*rr)["id"]

	Convey("Get user info by id: GET /user/id/:id", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/id/%v", api_v1, root_user_id))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, root_user_name)
		}
	})

	Convey("Update current user: PUT /user/update", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(map[string]string{
				"cnname": "cnroot2",
				"email":  "root2@test.com",
				"phone":  "18000000000",
			}).
			SetResult(rr).
			Put(fmt.Sprintf("%s/user/update", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "updated")
		}

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			client := resty.New()
			resp, err := client.R().
				SetHeader("X-Falcon-Token", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
			if err != nil {
				fmt.Println(err)
			} else {
				So(resp.StatusCode(), ShouldEqual, 200)
				So(*rr, ShouldNotBeEmpty)
				So((*rr)["cnname"], ShouldEqual, "cnroot2")
			}
		})
	})

	Convey("Change password: PUT /user/cgpasswd", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(map[string]string{
				"old_password": root_user_password,
				"new_password": root_user_password,
			}).
			SetResult(rr).
			Put(fmt.Sprintf("%s/user/cgpasswd", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "updated")
		}
	})

	Convey("Get user list: GET /user/users", t, func() {
		r := []map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(&r).
			Get(fmt.Sprintf("%s/user/users", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(r, ShouldNotBeEmpty)
			So(r[0]["name"], ShouldNotBeBlank)
		}
	})

	Convey("Get current user: POST /user/current", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/current", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, root_user_name)
		}
	})

	Convey("Auth user by session: GET /user/auth_session", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/auth_session", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "valid")
		}
	})

	Convey("Logout user: GET /user/logout", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/logout", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "successful")
		}
	})
}

func TestAdmin(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}
	var api_token string = ""

	Convey("Login as root", t, func() {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"name":     root_user_name,
				"password": root_user_password,
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, root_user_name)
			So((*rr)["sign"], ShouldNotBeBlank)
			So((*rr)["admin"], ShouldBeTrue)
		}
	})
	api_token = fmt.Sprintf(`{"name": "%v", "sign": "%v"}`, (*rr)["name"], (*rr)["sign"])
	Convey("Get user info by name: GET /user/name/:name", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
		}
	})
	test_user_id := (*rr)["id"]

	Convey("Change user role: PUT /admin/change_user_role", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"userID": %v,"admin": "yes"}`, test_user_id)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_role", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "changed")
		}

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			client := resty.New()
			resp, err := client.R().
				SetHeader("X-Falcon-Token", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
			if err != nil {
				fmt.Println(err)
			} else {
				So(resp.StatusCode(), ShouldEqual, 200)
				So(*rr, ShouldNotBeEmpty)
				So((*rr)["role"], ShouldEqual, 1)
			}
		})
	})

	Convey("Change user passwd: PUT /admin/change_user_passwd", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"userID": %v,"password": "%s"}`, test_user_id, test_user_password)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_passwd", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "updated")
		}
	})

	Convey("Change user profile: PUT /admin/change_user_profile", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"userID": %v,"cnname": "%s", "email": "%s"}`,
				test_user_id, test_user_name, "test_user1@test.com")).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_profile", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "updated")
		}

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			client := resty.New()
			resp, err := client.R().
				SetHeader("X-Falcon-Token", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
			if err != nil {
				fmt.Println(err)
			} else {
				So(resp.StatusCode(), ShouldEqual, 200)
				So(*rr, ShouldNotBeEmpty)
				So((*rr)["email"], ShouldEqual, "test_user1@test.com")
			}
		})
	})

	Convey("Admin login user: POST /admin/login", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(map[string]string{
				"name": test_user_name,
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/admin/login", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, test_user_name)
		}
	})

	Convey("Delete user: DELETE /admin/delete_user", t, func() {
	})
}

func TestTeam(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}

	Convey("Login as root", t, func() {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"name":     root_user_name,
				"password": root_user_password,
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, root_user_name)
			So((*rr)["sign"], ShouldNotBeBlank)
			So((*rr)["admin"], ShouldBeTrue)
		}
	})
	api_token := fmt.Sprintf(`{"name": "%v", "sign": "%v"}`, (*rr)["name"], (*rr)["sign"])

	Convey("Get user info by name: GET /user/name/:user", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["role"], ShouldEqual, 2)
			So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
		}
	})
	root_user_id := (*rr)["id"]

	Convey("Create team: POST /team", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(fmt.Sprintf(`{"name": "%s","resume": "i'm descript", "users": [%d]}`, test_team_name, test_user_id)).
			SetResult(rr).
			Post(fmt.Sprintf("%s/team", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "created")
		}
	})

	Convey("Get team by name: GET /team/name/:name", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().SetHeader("X-Falcon-Token", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/team/name/%s", api_v1, test_team_name))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, test_team_name)
			So((*rr)["users"], ShouldNotBeEmpty)
			So((*rr)["id"], ShouldBeGreaterThan, 0)
		}
	})
	test_team_id := (*rr)["id"]

	Convey("Get team by id: GET /team/id/:id", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().SetHeader("X-Falcon-Token", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/team/id/%v", api_v1, test_team_id))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["name"], ShouldEqual, test_team_name)
			So((*rr)["users"], ShouldNotBeEmpty)
			So((*rr)["id"], ShouldEqual, test_team_id)
		}
	})

	Convey("Update team by id: PUT /team", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(fmt.Sprintf(`{"id": %v,"resume": "descript2", "name":"%v", "users": [%d]}`,
				test_team_id, test_team_name, test_user_id)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/team", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "updated")
		}

		Convey("Get team by name: GET /team/name/:name", func() {
			*rr = map[string]interface{}{}
			client := resty.New()
			resp, err := client.R().SetHeader("X-Falcon-Token", api_token).SetResult(rr).
				Get(fmt.Sprintf("%s/team/name/%s", api_v1, test_team_name))
			if err != nil {
				fmt.Println(err)
			} else {
				So(resp.StatusCode(), ShouldEqual, 200)
				So(*rr, ShouldNotBeEmpty)
				So((*rr)["resume"], ShouldEqual, "descript2")
			}
		})
	})

	Convey("Add users to team: POST /team/user", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-Falcon-Token", api_token).
			SetBody(map[string]interface{}{
				"teamID": test_team_id,
				"users":  []string{root_user_name},
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/team/team/user", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "successful")
		}
	})

	Convey("Get teams which user belong to: GET /user/id/:id/teams", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().SetHeader("X-Falcon-Token", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/user/id/%v/teams", api_v1, root_user_id))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["teams"], ShouldNotBeEmpty)
		}
	})

	Convey("Check user in teams or not: GET /user/id/:id/in_teams", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetQueryParam("names", test_team_name).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/id/%v/in_teams", api_v1, root_user_id))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldEqual, "true")
		}
	})

	Convey("Get team list: GET /team", t, func() {
		var r []map[string]interface{}
		client := resty.New()
		resp, err := client.R().SetHeader("X-Falcon-Token", api_token).SetResult(&r).
			Get(fmt.Sprintf("%s/team", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(r, ShouldNotBeEmpty)
			So(r[0]["team"], ShouldNotBeEmpty)
			So(r[0]["users"], ShouldNotBeEmpty)
			So(r[0]["creator"], ShouldNotBeBlank)
		}
	})

	Convey("Delete team by id: DELETE /team/id/:id", t, func() {
		*rr = map[string]interface{}{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("X-Falcon-Token", api_token).
			SetResult(rr).
			Delete(fmt.Sprintf("%s/team/id/%v", api_v1, test_team_id))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["message"], ShouldContainSubstring, "deleted")
		}
	})
}

func TestGraph(t *testing.T) {
	api_token, err := get_session_token()
	if err != nil {
		log.Fatal(err)
	}

	client := resty.New()
	client.SetHeader("X-Falcon-Token", api_token)

	Convey("Get endpoint list: GET /graph/endpoint", t, func() {
		r := []map[string]interface{}{}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetQueryParams(map[string]string{
				"q":    ".+",
				"tags": "labels",
			}).
			SetResult(&r).
			Get(fmt.Sprintf("%s/graph/endpoint", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
			So(len(r), ShouldBeGreaterThanOrEqualTo, 0)
		}

		if len(r) == 0 {
			return
		}

		eid := r[0]["id"]
		r = []map[string]interface{}{}
		Convey("Get counter list: GET /graph/endpoint_counter", func() {
			resp, err := client.R().
				SetQueryParam("eid", fmt.Sprintf("%v", eid)).
				SetQueryParam("metricQuery", ".+").
				SetQueryParam("limit", "1").
				SetResult(&r).
				Get(fmt.Sprintf("%s/graph/endpoint_counter", api_v1))
			if err != nil {
				fmt.Println(err)
			} else {
				So(resp.StatusCode(), ShouldEqual, 200)
				So(r, ShouldNotBeEmpty)
			}
		})
	})
}

func TestNodata(t *testing.T) {
	api_token, err := get_session_token()
	if err != nil {
		log.Fatal(err)
	}

	var rr *map[string]interface{} = &map[string]interface{}{}
	client := resty.New()
	client.SetHeader("X-Falcon-Token", api_token)

	var nid int = 0

	Convey("Create nodata config: POST /nodata", t, func() {
		nodata_name := fmt.Sprintf("api.testnodata-%s", cu.RandString(8))
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"tags": "tags", "step": 60, "objType": "host", "obj": "docker-agent",
				"name": "%s", "mock": -1, "metric": "api.test.metric", "dsType": "GAUGE"}`, nodata_name)).
			SetResult(rr).
			Post(fmt.Sprintf("%s/nodata", api_v1))
		if err != nil {
			fmt.Println(err)
		} else {
			So(resp.StatusCode(), ShouldEqual, 200)
		}

		if v, ok := (*rr)["id"]; ok {
			nid = int(v.(float64))
			Convey("Delete nodata config", func() {
				resp, err := client.R().
					SetHeader("Content-Type", "application/json").
					Delete(fmt.Sprintf("%s/nodata/id/%d", api_v1, nid))
				if err != nil {
					fmt.Println(err)
				} else {
					So(resp.StatusCode(), ShouldEqual, 200)
				}
			})
		}
	})
}
