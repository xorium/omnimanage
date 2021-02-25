package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"omnimanage/internal/controller"
	omnimiddleware "omnimanage/internal/middleware"
	"omnimanage/internal/store"
	"omnimanage/internal/validator"
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
	// ...

	// logger
	// ...

	// init store
	db, err := getDB()
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
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = omniErr.ErrHandler

	// Middleware
	e.Use(
		omnimiddleware.ResponseType,
		middleware.Recover(),
		middleware.RequestID(),
	)

	//e.Use(middleware.Logger())

	// Controllers
	cntrManager := controller.NewManager(store)

	// Routes
	e.Use()
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

		// User relations
		userRoutes.GET("/:id/relationships/:rel", cntrManager.User.GetRelation)
		userRoutes.Match([]string{"PATCH", "POST", "DELETE"}, "/:id/relationships/:rel", cntrManager.User.ModifyRelation)
	}

	// Role routes
	{
		roleRoutes := companyGrp.Group("/roles")
		roleRoutes.GET("", cntrManager.Role.GetList)
		roleRoutes.GET("/:id", cntrManager.Role.GetOne)
	}
	// Start Server
	s := &http.Server{
		Addr:         ":8081",          // -> to config
		ReadTimeout:  10 * time.Minute, // -> to config
		WriteTimeout: 10 * time.Minute, // -> to config
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

func getDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		"5433",
		"db_user",
		"tYSk4dqaW7Hq4cw2r4hP",
		"omnimanage_db",
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
