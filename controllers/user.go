package controllers

import (
	"casbin-demo/models"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// UserController operations about Users
type UserController struct {
	beego.Controller
}

func init() {
	registerUserPolicy()
}

func registerUserPolicy() {
	// Path prefix, this is adjusted according to the specific project
	api := "/v1/user"
	// routing policy
	adminPolicy := map[string][]string{
		"/register": {"post"},
	}
	userPolicy := map[string][]string{
		// Note-use keyMatch2 in casbin.conf for obj
		// Verify, use: id to identify the parameter
		"/:id": {"get", "put", "delete"},
	}
	anonymousPolicy := map[string][]string{
		"/login": {"post"},
	}
	// models.RoleAdmin      = "admin"
	// models.RoleUser       = "user"
	// models.RoleAnonymous  = "anonymous"
	AddPolicyFromController(models.RoleAdmin, adminPolicy, api)
	AddPolicyFromController(models.RoleUser, userPolicy, api)
	AddPolicyFromController(models.RoleAnonymous, anonymousPolicy, api)
}

// AddPolicyFromController ...
func AddPolicyFromController(role string, policy map[string][]string, api string) {
	for path := range policy {
		for _, method := range policy[path] {
			// models.Enforcer is defined and initialized in models / Casbin.go
			_, _ = models.Enforcer.AddPolicy(models.GetRoleString(role), fmt.Sprintf("%s%s", api, path), method)
		}
	}
}

// @Title Register
// @Description Only administrator can register
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router  /register [post]
func (u *UserController) Register() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	// uid := models.AddUser(user)
	//u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJSON()
}

// @Title Profile
// @Description Only users and administrators can see the profiles of others or themselves. Because the administrator inherits the user, the administrator can do what the user can do
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router  /profile/:uid [get]
func (u *UserController) Profile() {
	// uid := u.GetString(":uid")
	// if uid != "" {
	// 	user, err := models.GetUser(uid)
	// 	if err != nil {
	// 		u.Data["json"] = err.Error()
	// 	} else {
	// 		u.Data["json"] = user
	// 	}
	// }
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	// username := u.GetString("username")
	// password := u.GetString("password")
	// if models.Login(username, password) {
	// 	u.Data["json"] = "login success"
	// } else {
	// 	u.Data["json"] = "user not exist"
	// }
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}
