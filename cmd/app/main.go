package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	"omnimanage/pkg/mapper"
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

	// ADMIN
	//{
	//	Admin := admin.New(&admin.AdminConfig{
	//		//DB: db.,
	//	})
	//
	//	Admin.AddResource(&src.User{})
	//	Admin.AddResource(&src.Role{})
	//	Admin.AddResource(&src.Location{})
	//	Admin.AddResource(&src.Company{})
	//
	//	adminHandler := echo.WrapHandler(Admin.NewServeMux("/admin"))
	//	e.GET("/admin", adminHandler)
	//	e.Any("/admin/*", adminHandler)
	//}

	// Middleware
	e.Use(
		omnimiddleware.ResponseType,
		middleware.Recover(),
		middleware.RequestID(),
	)

	//e.Use(middleware.Logger())

	// Controllers
	cntrManager := controller.NewManager(store, mapper.NewModelMapper())

	// Routes

	// Common grp
	companyGrp := e.Group("/companies/:idComp")

	// User routes
	{
		userRoutes := companyGrp.Group("/users")
		userRoutes.GET("", cntrManager.User.GetList)
		userRoutes.GET("/:id", cntrManager.User.GetOne)
		userRoutes.POST("/", cntrManager.User.Create)
		userRoutes.PATCH("/:id", cntrManager.User.Update)
		userRoutes.DELETE("/:id", cntrManager.User.Delete)

		// relations
		userRoutes.GET("/:id/relationships/:rel", cntrManager.User.GetRelation)
		userRoutes.Match(
			[]string{"PATCH", "POST", "DELETE"},
			"/:id/relationships/:rel",
			cntrManager.User.ModifyRelation)
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
