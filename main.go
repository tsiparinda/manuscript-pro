package main

import (
	"brucheion/handlers/api"
	"brucheion/handlers/newauth"
	"brucheion/models/repo"
	"brucheion/utils"
	"sync"

	// "github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http"
	// "github.com/vedicsociety/platform/logging"

	"brucheion/auth"
	"brucheion/handlers/admin"
	"brucheion/handlers/root"
	"brucheion/handlers/root/collectionedit"

	"github.com/vedicsociety/platform/authorization"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/pipeline"
	"github.com/vedicsociety/platform/pipeline/basic"
	"github.com/vedicsociety/platform/services"
	"github.com/vedicsociety/platform/sessions"
)

func registerServices() {
	services.RegisterDefaultServices()
	repo.RegisterSqlRepositoryService()
	sessions.RegisterSessionService()
	authorization.RegisterDefaultSignInService()
	authorization.RegisterDefaultUserService()
	auth.RegisterUserStoreService()
}

func createPipeline() pipeline.RequestPipeline {
	return pipeline.CreatePipeline(
		&basic.ServicesComponent{},
		&basic.LoggingComponent{},
		&basic.ErrorComponent{},
		&basic.AuthComponent{},
		&basic.StaticFileComponent{},
		&sessions.SessionComponent{},

		authorization.NewAuthComponent(
			"admin", // prefix
			authorization.NewRoleCondition("Administrators"), // condition
			admin.AdminHandler{},                             // requestHandlers ...interface{}
			admin.UsersHandler{},
			admin.GroupsHandler{},
			admin.RoutesHandler{},
			auth.SignOutHandler{},
		).AddFallback("/admin/section/", "^/admin[/]?$"),

		authorization.NewAuthComponent(
			"tools",
			authorization.NewRoleCondition("Users"),
			root.IngestCEXHandler{},
			root.AddColHandler{},
			root.IngestImageHandler{},
			root.ShareCollectionHandler{},
			collectionedit.CollectionEditHandler{},
			root.CollectionBulkEditHandler{},
			auth.AccountHandler{},
			root.MenuHandler{},
		).AddFallback("/tools/", "^[/]?$"),

		authorization.NewAuthComponent(
			"api/v2",
			authorization.NewRoleCondition("Users"),
			// 	api.PassageHandler{},
			// 	api.CollectionHandler{},
			api.UsersHandler{},
			api.CEXuploadHandler{},
			api.AddColHandler{},
			api.ShareCollectionHandler{},
			api.EditPassageHandler{},
			api.DropCollectionHandler{},
			api.IngestImageHandler{},
		).AddFallback("/api/v2/", "^/api[/]?$"),

		handling.NewRouter(
			//	handling.HandlerEntry{"", collections.CollectionsHandler{}},
			handling.HandlerEntry{"", root.CollectionsHandler{}},
			handling.HandlerEntry{"", root.CollectionOverviewHandler{}},
			// handling.HandlerEntry{"", root.IngestImageHandler{}},
			handling.HandlerEntry{"", root.MenuHandler{}},
			// handling.HandlerEntry{"", root.IngestCEXHandler{}},
			handling.HandlerEntry{"", auth.AuthenticationHandler{}},
			handling.HandlerEntry{"auth", newauth.NewAuthenticationHandler{}},
			handling.HandlerEntry{"api/v1", api.CollectionsHandler{}},
			handling.HandlerEntry{"api/v1", api.CollectionHandler{}},
			handling.HandlerEntry{"api/v1", api.PassageHandler{}},
			handling.HandlerEntry{"api/v1", api.UserHandler{}},
			handling.HandlerEntry{"api/v1", api.ImagesHandler{}},
			// handling.HandlerEntry{"api/v1", api.UserHandler{}},
			// handling.HandlerEntry{"api/v1", api.CEXuploadHandler{}},
			//		).AddMethodAlias("/", tools.ToolsHandler.GetSection, ""))
		).AddMethodAlias("/", root.CollectionsHandler.GetCollections, 0, 1).
			AddMethodAlias("/collections[/]?[A-z0-9]*?",
				root.CollectionsHandler.GetCollections, 0, 1))
}

func main() {

	registerServices()
	//repo := services.GetService(repo.SqlRepository)
	// go utils.ImageHelper.SyncDZIcollection(
	// 	utils.ImageHelper{
	// 		Repository: repo.SqlRepository,
	// 		Logger:        logging.Logger,
	// 		Config:        config.Configuration,
	// 	})
	go services.Call(utils.SyncDZIcollection)
	//services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	services.Call(utils.InitProviders) // init oauth goth providers
	services.Call(utils.SeedAdmin)     // seed administartor user to DB
	//repo.LoadMigrations(config.Configuration, logging.Logger)

	results, err := services.Call(http.Serve, createPipeline())
	if err == nil {
		(results[0].(*sync.WaitGroup)).Wait()
	} else {
		panic(err)
	}

}
