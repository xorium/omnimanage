package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"omnimanage/internal/config"
	"omnimanage/internal/controller"
	omnimiddleware "omnimanage/internal/middleware"
	"omnimanage/internal/store"
	omniErr "omnimanage/pkg/error"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	//ctx := context.Background()

	// config
	cfg, err := config.Get("")
	if err != nil {
		log.Fatal("can't read config: %v", err)
	}

	// logger
	// ...

	// init store
	db, err := getDB(cfg)
	if err != nil {
		return err
	}

	store := store.NewStore(db)

	// Init service manager
	//serviceManager, err := service.NewManager(store)
	//if err != nil {
	//	return fmt.Errorf("store.New failed: %w", err)
	//}

	// Init echo instance
	e := echo.New()
	e.Debug = cfg.App.Debug
	//e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = omniErr.ErrHandler

	// Middleware
	e.Use(
		omnimiddleware.ResponseType,
		middleware.Recover(),
		middleware.RequestID(),
		middleware.RemoveTrailingSlash(),
	)

	//e.Use(middleware.Logger())

	// Controllers
	cntrManager := controller.NewManager(store)

	//// Docs
	se := echoswagger.New(e, "docs/", &echoswagger.Info{
		Title:          "Swagger Omnimanage",
		Description:    "Omnimanage description.  You can find out more about     Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).      For this sample, you can use the api key `special-key` to test the authorization     filters.",
		Version:        "1.0.0",
		TermsOfService: "http://swagger.io/terms/",
		Contact: &echoswagger.Contact{
			Email: "apiteam@swagger.io",
		},
		License: &echoswagger.License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
	})
	se.SetExternalDocs("Find out more about Swagger", "http://swagger.io").
		SetResponseContentType("application/xml", "application/json").
		SetUI(echoswagger.UISetting{DetachSpec: true, HideTop: true}).
		SetScheme("https", "http")

	// Routes

	// Company
	companyGrp := e.Group("/companies/:idComp") //se.Group("Company", "/companies")
	//{
	//	companyGrp.GET("", nil)
	//	companyGrp.GET("/:cid", nil).AddParamPath(0, "cid", "Company ID")
	//}

	// User routes
	{

		userRoutes := companyGrp.Group("/users")
		//userRoutes := se.Group("Users", "/company/users")
		userRoutes.GET("", cntrManager.User.GetList)
		userRoutes.GET("/:id", cntrManager.User.GetOne)
		userRoutes.POST("/", cntrManager.User.Create)
		userRoutes.PATCH("/:id", cntrManager.User.Update)
		userRoutes.DELETE("/:id", cntrManager.User.Delete)

		// relations
		userRoutes.GET("/:id/relationships/:rel", cntrManager.User.GetRelation)
		//userRoutes.Match(
		//	[]string{"PATCH", "POST", "DELETE"},
		//	"/:id/relationships/:rel",
		//	cntrManager.User.ModifyRelation)
	}

	// Role routes
	{
		roleRoutes := companyGrp.Group("/roles")
		roleRoutes.GET("", cntrManager.Role.GetList)
		roleRoutes.GET("/:id", cntrManager.Role.GetOne)
	}

	// Start Server
	s := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		ReadTimeout:  time.Millisecond * time.Duration(cfg.Server.ReadTimeoutMSec),
		WriteTimeout: time.Millisecond * time.Duration(cfg.Server.ReadTimeoutMSec),
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Printf("shutting down the server %v", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

func getDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	gormConf := &gorm.Config{}
	if cfg.App.Debug {
		gormConf.Logger = logger.Default.LogMode(logger.Info)
	}
	return gorm.Open(postgres.Open(dsn), gormConf)
}
