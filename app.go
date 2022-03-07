package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	admin_controller "github.com/asishshaji/admin-api/controller"
	"github.com/asishshaji/admin-api/utils"

	"github.com/asishshaji/admin-api/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	app  *echo.Echo
	port string
}

type Controllers struct {
	AdminController admin_controller.IAdminController
}

func NewApp(port string, controller Controllers) *App {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Secure())

	e.POST("/login", controller.AdminController.Login)
	adminGroup := e.Group("/admin")

	adminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &models.AdminJWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	adminGroup.Use(utils.AdminAuthenticationMiddleware)

	adminGroup.POST("/task", controller.AdminController.CreateTask)
	adminGroup.PUT("/task", controller.AdminController.UpdateTask)
	adminGroup.GET("/task", controller.AdminController.GetTasks)
	adminGroup.DELETE("/task", controller.AdminController.DeleteTask)

	adminGroup.GET("/users", controller.AdminController.GetUsers)

	adminGroup.GET("/submission", controller.AdminController.GetTaskSubmissions)
	adminGroup.PUT("/submission", controller.AdminController.EditTaskSubmissionStatus)

	adminGroup.GET("/user/submission/:id", controller.AdminController.GetTaskSubmissionForUser)

	adminGroup.GET("/mentor", controller.AdminController.GetMentors)
	adminGroup.POST("/mentor", controller.AdminController.CreateMentor)
	adminGroup.PUT("/mentor", controller.AdminController.UpdateMentor)

	return &App{
		app:  e,
		port: port,
	}
}

func (a *App) RunServer() {

	go func() {
		a.app.Logger.Fatal(a.app.Start(a.port))
	}()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sigData := <-sigChan

	log.Printf("Signal received : %v\n", sigData)
	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	a.app.Shutdown(tc)
}
