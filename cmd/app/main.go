package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"omnimanage/internal/controller"
	"omnimanage/internal/middleware"
	"omnimanage/internal/store"
	"omnimanage/internal/validator"
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
	//e.HTTPErrorHandler =

	// Middleware
	e.Use(middleware.ResponseType)
	//e.Use(middleware.Logger()) e.Use(middleware.Recover())

	// Controllers
	cntrManager := controller.NewManager(store)

	// Routes
	e.Use()
	// Common grp
	companyGrp := e.Group("/companies/:idComp")

	// User routes
	userRoutes := companyGrp.Group("/users")
	userRoutes.GET("", cntrManager.User.GetList)
	userRoutes.GET("/:id", cntrManager.User.GetOne)
	userRoutes.GET("/:id/relationships/:rel", cntrManager.User.GetRelation)

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

	//graceful.ListenAndServe()
	//us, err := store.Users.GetById(ctx, 1413)
	//if err != nil {
	//	return err
	//}
	//us.FirstName = "AAA"
	//db.Debug().WithContext(ctx).Save(&us)
	//
	//fmt.Printf("user: %v", us)

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
