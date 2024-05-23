package admin

import (
	"brucheion/models"
	"brucheion/utils"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	sessions "github.com/vedicsociety/platform/sessions"
)

type GroupsHandler struct {
	models.Repository
	handling.URLGenerator
	sessions.Session
	logging.Logger
}

type GroupsTemplateContext struct {
	Groups  []models.Group
	EditId  int
	EditUrl string
	SaveUrl string
}

const GROUP_EDIT_KEY string = "group_edit"

func (handler GroupsHandler) GetData() actionresults.ActionResult {
	return actionresults.NewTemplateAction("admin_groups.html",
		GroupsTemplateContext{
			Groups: handler.Repository.GetGroups(),
			EditId: handler.Session.GetValueDefault(GROUP_EDIT_KEY, 0).(int),
			EditUrl: utils.MustGenerateUrl(handler.URLGenerator,
				GroupsHandler.PostGroupEdit),
			SaveUrl: utils.MustGenerateUrl(handler.URLGenerator,
				GroupsHandler.PostGroupSave),
		})
}

type GroupEditReference struct {
	ID int
}

func (handler GroupsHandler) PostGroupEdit(ref EditReference) actionresults.ActionResult {
	handler.Session.SetValue(GROUP_EDIT_KEY, ref.ID)
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Groups"))
}

type GroupSaveReference struct {
	Id   int
	Name string
}

func (handler GroupsHandler) PostGroupSave(p GroupSaveReference) actionresults.ActionResult {
	handler.Logger.Debugf("PostGroupSave group: ", p)
	group := handler.Repository.GetGroupById(p.Id)
	if group.Id == 0 {
		handler.Logger.Panicf("PostGroupSave Group not found: ", p)
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AdminHandler.GetSection, "Groups"))
	}
	group.Name = p.Name
	handler.Repository.UpdateGroup(&group)
	handler.Session.SetValue(GROUP_EDIT_KEY, 0)
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Groups"))
}

func (handler GroupsHandler) GetSelect(r []models.Group) actionresults.ActionResult {
	// r - slice user's groups
	// we need to filter out this groups
	allg := handler.GetGroups()
	var fg []models.Group

	// looking for string in slice
	contains := func(g models.Group, s []models.Group) bool {
		for i := range s {
			if s[i].Name == g.Name {
				return true
			}
		}
		return false
	}
	// append string ito new slice
	for g := range allg {
		if !contains(allg[g], r) {
			fg = append(fg, allg[g])
		}
	}
	handler.Logger.Debugf("GetSelect Group in: ", r)
	return actionresults.NewTemplateAction("admin_select_group.html", struct {
		Current int
		Groups  []models.Group
	}{Current: 0, Groups: fg})
}
