package admin

import (
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	sessions "github.com/vedicsociety/platform/sessions"
)

// var sectionNames = []string{"Users", "Groups", "Database"}
var sectionNames = []string{"Users", "Groups", "Routes"}

type AdminHandler struct {
	handling.URLGenerator
	logging.Logger
	sessions.Session
}

type AdminTemplateContext struct {
	Sections       []string
	ActiveSection  string
	SectionUrlFunc func(string) string
}

func (handler AdminHandler) GetSection(section string) actionresults.ActionResult {
	handler.Logger.Debugf("GetSection section ", section)
	if section == "" {
		handler.Session.SetValue(USER_EDIT_KEY, 0)
		handler.Session.SetValue(GROUP_EDIT_KEY, 0)
	}
	return actionresults.NewTemplateAction("admin.html", AdminTemplateContext{
		Sections:      sectionNames,
		ActiveSection: section,
		SectionUrlFunc: func(sec string) string {
			sectionUrl, _ := handler.GenerateUrl(AdminHandler.GetSection, sec)
			return sectionUrl
		},
	})
}

