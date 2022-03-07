package main

import (
	"context"
	"log"
	"os"

	admin_controller "github.com/asishshaji/admin-api/controller"
	admin_repository "github.com/asishshaji/admin-api/repositories"
	"github.com/asishshaji/admin-api/services/admin_service"
	"github.com/asishshaji/admin-api/services/image_service"
	"github.com/asishshaji/admin-api/utils"
	"github.com/go-redis/redis/v8"
)

func main() {

	logger := log.New(os.Stdout, "admin-api", log.LstdFlags)

	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		logger.Println(err)
	} else {
		logger.Println("Connected to redis")
	}

	imageService := image_service.NewImageService(logger)

	adminRepo := admin_repository.NewAdminRepository(logger, db)
	adminService := admin_service.NewAdminService(logger, adminRepo, redisClient, imageService)
	adminController := admin_controller.NewAdminController(logger, adminService)

	controller := Controllers{
		AdminController: adminController,
	}

	app := NewApp(env.ServerPort, controller)
	app.RunServer()
}
