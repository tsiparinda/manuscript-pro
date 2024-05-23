// Package root provides API handlers for managing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// MenuHandler is a struct that embeds the necessary dependencies for handling menu-related requests.
// It implements methods that handle HTTP requests related to this process.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// an identity.User to access the user's identity, and a Logger from the logging package to log the operations
type MenuHandler struct {
	Repository   models.Repository
	URLGenerator handling.URLGenerator
	identity.User
	logging.Logger
}

// MenuTemplateContext is a struct that defines the context for rendering a menu template.
// It contains user details and a list of menu sections, among other properties.
type MenuTemplateContext struct {
	UserName       string
	MenuSections   []MainMenu
	ActiveSection  string
	User           *identity.User
	SectionUrlFunc func(string, string, string) string
}

// MenuItem defines a single item in a menu.
type MenuItem struct {
	Handler string
	Action  string
	Section string
	Name    string
}

// MainMenu represents a main menu, which is composed of multiple MenuItems.
type MainMenu struct {
	Paragraph string
	Items     []MenuItem
}

// baseMenu represents the default menu structure.
var baseMenu = []MainMenu{
	{"Tools",
		[]MenuItem{ //{Handler: "Tools", Action: "GETSECTION", Section: "PassageOverview", Name: "Passage overview"},
			{Handler: "IngestCEX", Action: "GETINGESTCEX", Section: "", Name: "Ingest CEX"},
			// {Handler: "CollectionBulkEdit", Action: "GETBULKEDIT", Section: "0", Name: "Bulk edit"},
		},
	},
	{"Account",
		[]MenuItem{{Handler: "Account", Action: "GETACCOUNTPROFILE", Section: "", Name: "Your profile"},
			{Handler: "Authentication", Action: "GETSIGNOUT", Section: "", Name: "Sign out"}},
	},
}

// GetMenu is a method on MenuHandler that handles a GET request to retrieve the menu.
// It logs the operation and then returns an ActionResult that represents a server-side rendered
// template response. The template file used is 'root_menu.html', and it receives a MenuTemplateContext
// instance that contains the necessary data for rendering the menu.
func (h MenuHandler) GetMenu(section string) actionresults.ActionResult {
	// userkey := handler.Session.GetValueDefault(USER_SESSION_KEY, 0).(int)
	h.Logger.Debugf("GetMenu ", section, h.User, h.User.InRole("Administrators"))

	// add admin menu item
	var newMenu []MainMenu
	if h.User.InRole("Administrators") {
		for _, p := range baseMenu {
			if p.Paragraph == "Account" {
				var items []MenuItem
				items = append(items, MenuItem{Handler: "Admin", Action: "GETSECTION", Section: "", Name: "Admin panel"})
				items = append(items, p.Items...)
				p.Items = []MenuItem{}
				p.Items = make([]MenuItem, len(items))
				copy(p.Items, items)
				newMenu = append(newMenu, p)
			} else {
				newMenu = append(newMenu, p)
			}
		}
	} else {
		newMenu = append(newMenu, baseMenu...)
	}
	// change Account item menu as user's name
	for p := range newMenu {
		if newMenu[p].Paragraph == "Account" {
			newMenu[p].Paragraph = h.GetDisplayName()
		}
	}

	return actionresults.NewTemplateAction("root_menu.html",
		MenuTemplateContext{
			UserName:      h.User.GetDisplayName(), //( handler.Session.GetValueDefault(USER_SESSION_KEY, 0).(int)),
			MenuSections:  newMenu,
			ActiveSection: section,
			User:          &h.User,
			SectionUrlFunc: func(hand, act, sec string) string {
				sectionUrl, _ := h.URLGenerator.GenerateURLByName(hand, act, sec)
				return sectionUrl
			},
		})
}
