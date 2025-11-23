package main

import (
	"fmt"
	"log"
	"onboarding/api"
	"onboarding/internal/handler"
	"onboarding/internal/repository"
	"onboarding/internal/service"
	"onboarding/pkg/config"
	"onboarding/pkg/database"
	"onboarding/pkg/token"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Init DB error: %v", err)
	}

	jwtImpl, err := token.NewJWT(cfg.Token)
	if err != nil {
		log.Fatalf("Couldn't create token maker: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, jwtImpl)
	authHandler := handler.NewAuthHandler(authService)

	server := api.NewServer(cfg.App, jwtImpl, authHandler, userHandler)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	err = server.Start(fmt.Sprintf(":%s", cfg.App.Port))
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}
