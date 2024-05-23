package admin

import (
	"brucheion/models"
	"fmt"
	"reflect"
	"regexp"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	sessions "github.com/vedicsociety/platform/sessions"
)

// handler for showing routes and its parameters
// useful for debugging
type RoutesHandler struct {
	models.Repository
	handling.URLGenerator
	sessions.Session
	logging.Logger
}

type Route struct {
	HttpMethod    string
	Prefix        string
	HandlerName   string
	ActionName    string
	Expression    regexp.Regexp
	HandlerMethod reflect.Method
}

type RoutesTemplateContext struct {
	Routes []string
}

func (handler RoutesHandler) GetData() actionresults.ActionResult {
	var routes []string
	for _, r := range handler.URLGenerator.RoutesPrint() {
		routes = append(routes, fmt.Sprintln("", r))
	}
	// testurl, _ := handler.GenerateURLByName("Admin", "GETSIGNOUT", "SignOut")
	// handler.Logger.Debugf("GetSection Route", testurl)
	// []handling.Route
	return actionresults.NewTemplateAction("admin_routes.html",
		RoutesTemplateContext{
			Routes: routes,
		})
}
