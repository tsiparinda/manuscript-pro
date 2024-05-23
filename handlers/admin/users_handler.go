package admin

import (
	"brucheion/models"
	"brucheion/utils"
	"encoding/json"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	sessions "github.com/vedicsociety/platform/sessions"
)

type UsersHandler struct {
	models.Repository
	handling.URLGenerator
	sessions.Session
	logging.Logger
	config.Configuration
}

type UsersTemplateContext struct {
	Users        []models.Credentials
	EditId       int
	EditUrl      string
	SaveUrl      string
	CancelUrl    string
	DeleteUrl    string
	ForgotPwdUrl string
}

const USER_EDIT_KEY string = "user_edit"

func (handler UsersHandler) GetData() actionresults.ActionResult {
	//handler.Logger.Debugf("GetUsers", handler.Repository.GetUsers())
	return actionresults.NewTemplateAction("admin_users.html",
		UsersTemplateContext{
			Users:  handler.Repository.GetCredentials(),
			EditId: handler.Session.GetValueDefault(USER_EDIT_KEY, 0).(int),
			EditUrl: utils.MustGenerateUrl(handler.URLGenerator,
				UsersHandler.PostUserEdit),
			SaveUrl: utils.MustGenerateUrl(handler.URLGenerator,
				UsersHandler.PostUserSave),
			CancelUrl: utils.MustGenerateUrl(handler.URLGenerator,
				AdminHandler.GetSection, ""),
			DeleteUrl: utils.MustGenerateUrl(handler.URLGenerator, UsersHandler.PostUserDelete),
			//ForgotPwdUrl: utils.MustGenerateUrl(handler.URLGenerator, UsersHandler.PostUserResetPassword),
		})
}

type EditReference struct {
	ID int
}

func (handler UsersHandler) PostUserEdit(ref EditReference) actionresults.ActionResult {
	handler.Session.SetValue(USER_EDIT_KEY, ref.ID)
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Users"))
}

type UserSaveReference struct {
	Id       int
	Username string
	Email    string
	Action   string
	Roles    string
}

func (handler UsersHandler) PostUserSave(p UserSaveReference) actionresults.ActionResult {
	handler.Logger.Debugf("PostUserSave User: ", p)

	switch p.Action {
	case "Save":
		user := handler.Repository.GetUserByID(p.Id)
		if user.Id == 0 {
			handler.Logger.Panicf("PostUserSave User not found: ", p)
			return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
				AdminHandler.GetSection, "Users"))
		}
		user.Email = p.Email
		user.Username = p.Username

		var rolesid []int
		json.Unmarshal([]byte(p.Roles), &rolesid)
		if p.Id != 1 {
			handler.Logger.Debugf("PostUserSave rolesid: ", rolesid)
			//roleslice := strings.Split(p.Roles, "\n")
			user.Roles = []models.Group{}
			for r := range rolesid {
				user.Roles = append(user.Roles, handler.Repository.GetGroupById(rolesid[r]))
			}
			if err := handler.Repository.UpdateUserGroups(&user); err != nil {
				handler.Logger.Panicf("PostUserSave Cannot exec UpdateUserGroups command: %v", err.Error())
			}
		} else {
			handler.Logger.Debugf("PostUserSave: you can't delete the Administrator account!")
		}
		handler.Logger.Debugf("PostUserSave User: ", user)
		if err := handler.Repository.UpdateUser(&user); err != nil {
			handler.Logger.Panicf("PostUserSave Cannot exec UpdateUser command: %v", err.Error())
		}
	case "ResetPwd":
		handler.Logger.Debugf("PostUserSave resetpwd User: ", p)
		user := handler.Repository.GetUserByID(p.Id)
		if user.Id == 0 {
			handler.Logger.Panicf("PostUserSave User not found: ", p)
			return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
				AdminHandler.GetSection, "Users"))
		}
		user.Email = p.Email
		user.Username = p.Username
		err := utils.ResetPassword(handler.Repository, handler.Session, handler.URLGenerator, handler.Logger, handler.Configuration, user.Email)
		if err != nil {
			handler.Logger.Debugf("PostForgotPwd An error arise after call ResetPassword: ", err.Error())
			resp := models.JSONResponse{
				Status:  "success",
				Message: "The error arised! " + err.Error(),
			}
			return actionresults.NewErrorJsonAction(resp)
		}
	}
	handler.Session.SetValue(USER_EDIT_KEY, 0)
	handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "")

	resp := models.JSONResponse{
		Status:  "success",
		Message: "Link for reset password was sent to " + p.Email,
	}
	return actionresults.NewJsonAction(resp)

}

func (handler UsersHandler) PostUserDelete(p UserSaveReference) actionresults.ActionResult {
	handler.Logger.Debugf("PostUserDelete User: ", p.Id)
	// here should be code for delete user
	handler.Session.SetValue(USER_EDIT_KEY, 0)
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Users"))
}
