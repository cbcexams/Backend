package controllers

import (
	"cbc-backend/models"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description user login
// @Param	body	body	models.User	true	"Username and password"
// @Success 200 {string} token
// @Failure 403 user not exist
// @router /login [post]
func (u *UserController) Login() {
	var loginInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &loginInfo)
	if err != nil {
		u.Data["json"] = map[string]string{"error": "Invalid request body"}
		u.ServeJSON()
		return
	}

	user, err := models.GetUser(loginInfo.Username)
	if err != nil {
		u.Data["json"] = map[string]string{"error": "User not found"}
		u.ServeJSON()
		return
	}

	if err := user.ValidatePassword(loginInfo.Password); err != nil {
		u.Data["json"] = map[string]string{"error": "Invalid password"}
		u.ServeJSON()
		return
	}

	token, err := user.GenerateToken()
	if err != nil {
		u.Data["json"] = map[string]string{"error": err.Error()}
		u.ServeJSON()
		return
	}

	// Get session ID from request
	sessionID := u.Ctx.Input.Cookie("session_id")
	if sessionID == "" {
		sessionID = u.Ctx.Input.Header("X-Session-ID")
	}

	// Link session to user if session exists
	if sessionID != "" {
		err = models.LinkSessionToUser(sessionID, user.Id)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to link session to user: %v\n", err)
		}
	}

	u.Data["json"] = map[string]interface{}{
		"token":      token,
		"session_id": sessionID,
	}
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

// @Title Signup
// @Description create new user
// @Param	body	body	models.User	true	"User info"
// @Success 200 {string} token
// @Failure 403 body is empty
// @router /signup [post]
func (u *UserController) Signup() {
	var user models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		u.Data["json"] = map[string]string{"error": "Invalid request body"}
		u.ServeJSON()
		return
	}

	err = models.AddUser(&user)
	if err != nil {
		u.Data["json"] = map[string]string{"error": err.Error()}
		u.ServeJSON()
		return
	}

	token, err := user.GenerateToken()
	if err != nil {
		u.Data["json"] = map[string]string{"error": err.Error()}
		u.ServeJSON()
		return
	}

	// Get session ID from request
	sessionID := u.Ctx.Input.Cookie("session_id")
	if sessionID == "" {
		sessionID = u.Ctx.Input.Header("X-Session-ID")
	}

	// Link session to user if session exists
	if sessionID != "" {
		err = models.LinkSessionToUser(sessionID, user.Id)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to link session to user: %v\n", err)
		}
	}

	u.Data["json"] = map[string]interface{}{
		"token":      token,
		"session_id": sessionID,
	}
	u.ServeJSON()
}
