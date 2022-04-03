package main

import (
	"context"
	"log"
	"os"

	admin_controller "github.com/asishshaji/admin-api/controller"
	admin_repository "github.com/asishshaji/admin-api/repositories"
	"github.com/asishshaji/admin-api/services/admin_service"
	file_service "github.com/asishshaji/admin-api/services/file"
	"github.com/asishshaji/admin-api/services/notification_service"
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

	fileService := file_service.NewFileService(logger)
	onesignalService := notification_service.NewNotificationService(logger)

	adminRepo := admin_repository.NewAdminRepository(logger, db)
	adminService := admin_service.NewAdminService(logger, adminRepo, redisClient, onesignalService)
	adminController := admin_controller.NewAdminController(logger, adminService, fileService)

	password, err := utils.Hashpassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Fatalln("Error creating admin")
	}

	err = adminRepo.GenerateAdminCredentials(context.Background(), os.Getenv("ADMIN_USERNAME"), password)
	if err != nil {
		log.Fatalln("Error creating admin")
	}

	controller := Controllers{
		AdminController: adminController,
	}

	app := NewApp(env.ServerPort, controller)
	app.RunServer()
}
